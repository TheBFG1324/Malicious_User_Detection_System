"""
Uploaded Interactions and Associations to the Neo4j database

This script:
 1) Connects to a Neo4j database.
 2) Creates up to 10 fake users (by default) with a random malicious_score in [0,1].
 3) Generates 0-100 interactions per user, with more interactions (and more honeytoken triggers)
    for users who have higher maliciousness.
 4) Creates ASSOCIATED_WITH relationships among users, with higher-maliciousness
    users more likely to associate with each other (forming a "bad cluster").

Helped generated by OpenAI's o1 model. 
"""

import random
from neo4j import GraphDatabase

class DataUploader:
    def __init__(self, uri, user, password):
        self.driver = GraphDatabase.driver(uri, auth=(user, password))

    def close(self):
        self.driver.close()

    def upload_data(self, num_users=10):
        """
        Generates fake data for 'num_users' and uploads it to Neo4j.
        """
        users = self._create_fake_users(num_users)
        self._create_interactions(users)
        self._create_associations(users)
        print("Data upload complete!")

    def _create_fake_users(self, num_users):
        """
        Creates 'num_users' in-memory user objects (each with a user_id and malicious_score),
        then MERGEs them into Neo4j.
        """
        users = []
        with self.driver.session() as session:
            for i in range(num_users):
                user_id = f"user_{i}"
                malicious_score = round(random.random(), 3)  # in [0,1]
                
                users.append({"user_id": user_id, "malicious_score": malicious_score})
                
                session.execute_write(
                    self._merge_user, user_id, malicious_score
                )
        return users

    @staticmethod
    def _merge_user(tx, user_id, malicious_score):
        """
        Merges a User node with a given malicious_score. If it doesn't exist,
        creates it with malicious_score=..., otherwise updates it.
        """
        query = """
        MERGE (u:User {user_id: $user_id})
        ON CREATE SET u.malicious_score = $malicious_score
        ON MATCH SET u.malicious_score = $malicious_score
        """
        tx.run(query, user_id=user_id, malicious_score=malicious_score)

    def _create_interactions(self, users):
        """
        Creates a set of interactions for each user. The number and nature
        of interactions roughly correlates with the user's malicious_score:
         - 0 to 100 total interactions
         - More malicious => more interactions
         - Higher ratio of honeytoken triggers for higher malicious_score
        """
        with self.driver.session() as session:
            for user in users:
                score = user["malicious_score"]

                # Base random from 0..50, then add up to 50 more if score is high
                base_interactions = random.randint(0, 50)
                extra_interactions = int(50 * score)  # up to 50
                total_interactions = base_interactions + extra_interactions
                if total_interactions == 0:
                    # Ensure at least 1 interaction for demonstration
                    total_interactions = 1

                # Honeytoken triggers as fraction of total, weighted by malicious_score
                # e.g. if score is high, we trigger more honeytoken
                # random.uniform(0.1, 0.9) to add variety
                honeytoken_fraction = random.uniform(0.1, 0.9) * score
                honeytoken_count = int(total_interactions * honeytoken_fraction)

                # Generate a small pool of IP addresses
                # More malicious => more IP addresses
                ip_pool_size = random.randint(1, 3) + int(4 * score)  # up to ~7 IPs
                ip_pool = [
                    f"{random.randint(1,255)}.{random.randint(1,255)}."
                    f"{random.randint(1,255)}.{random.randint(1,255)}"
                    for _ in range(ip_pool_size)
                ]

                for _ in range(total_interactions):
                    endpoint = f"/api/endpoint_{random.randint(1,5)}"
                    response_status_code = random.choice([200, 400, 401, 404, 500])
                    
                    if honeytoken_count > 0:
                        honeytoken_triggered = True
                        honeytoken_count -= 1
                    else:
                        honeytoken_triggered = False
                    
                    ip_address = random.choice(ip_pool)

                    interaction_params = {
                        "user_id": user["user_id"],
                        "endpoint": endpoint,
                        "response_status_code": response_status_code,
                        "honeytoken_triggered": honeytoken_triggered,
                        "ip_address": ip_address
                    }

                    session.execute_write(
                        self._create_interaction, interaction_params
                    )

    @staticmethod
    def _create_interaction(tx, interaction_params):
        """
        Inserts an Interaction node connected to the user with :HAS_INTERACTION.
        """
        query = """
        MERGE (u:User {user_id: $user_id})
            ON CREATE SET u.malicious_score = 0.0
        CREATE (i:Interaction {
            endpoint: $endpoint,
            timestamp: timestamp(),
            response_status_code: $response_status_code,
            honeytoken_triggered: $honeytoken_triggered,
            ip_address: $ip_address
        })
        CREATE (u)-[:HAS_INTERACTION]->(i)
        """
        tx.run(query, **interaction_params)

    def _create_associations(self, users):
        """
        Creates ASSOCIATED_WITH relationships among users. Higher malicious
        users are more likely to associate with each other.
        """
        with self.driver.session() as session:
            # Sort users by malicious_score ascending
            sorted_users = sorted(users, key=lambda u: u["malicious_score"])
            
            for i in range(len(sorted_users)):
                user = sorted_users[i]
                user_id = user["user_id"]
                user_score = user["malicious_score"]

                # Decide how many associations to create
                # We'll do 1..3 associations per user
                # Higher malicious => more associations, typically to other malicious users
                assoc_count = random.randint(1, 3) + int(user_score * 2)

                # If user is high malicious (score>0.5), pick from top half
                if user_score > 0.5:
                    start_idx = len(sorted_users)//2
                    end_idx = len(sorted_users) - 1
                else:
                    # If user is lower malicious, pick from bottom half
                    start_idx = 0
                    end_idx = len(sorted_users)//2

                candidates = [u for u in sorted_users[start_idx:end_idx+1] if u["user_id"] != user_id]
                # If too few in that half, just pick from the entire set
                if len(candidates) < assoc_count:
                    candidates = [u for u in sorted_users if u["user_id"] != user_id]

                # Choose 'assoc_count' random users from the candidate list
                chosen = random.sample(candidates, k=min(assoc_count, len(candidates)))

                for assoc_user in chosen:
                    session.execute_write(self._associate_users, user_id, assoc_user["user_id"])

    @staticmethod
    def _associate_users(tx, user1, user2):
        """
        Creates a bi-directional :ASSOCIATED_WITH relationship between two users.
        """
        query = """
        MERGE (u1:User {user_id: $user1})
            ON CREATE SET u1.malicious_score = 0.0
        MERGE (u2:User {user_id: $user2})
            ON CREATE SET u2.malicious_score = 0.0
        CREATE (u1)-[:ASSOCIATED_WITH]->(u2)
        CREATE (u2)-[:ASSOCIATED_WITH]->(u1)
        """
        tx.run(query, user1=user1, user2=user2)


def main():
    uri = "bolt://localhost:7687"
    username = "neo4j"
    password = "Password"

    uploader = DataUploader(uri, username, password)
    try:
        uploader.upload_data(num_users=10) 
    finally:
        uploader.close()

if __name__ == "__main__":
    main()
