# Creates a synthetic dataset with to train AI model on maliciousness score prediction
# Helped generated by OpenAI's o1 model. 

# Required Libraries
import csv
import random

# GenerateData class to generate synthetic data
class GenerateData:
    def __init__(self, filename):
        self.filename = filename

    def generate_data(self):
        with open(self.filename, mode='w', newline='') as data_file:
            data_writer = csv.writer(data_file, delimiter=',', quotechar='"', quoting=csv.QUOTE_MINIMAL)
            # Write CSV header
            data_writer.writerow([
                'total_access_count',
                'honeytoken_access_count',
                'shared_ip_count',
                'avg_associated_malicious_score',
                'actual_maliciousness_score'
            ])

            # Generate 1000 rows of synthetic data
            for _ in range(1000):
                # 1. Generate correlated features
                total_access_count = random.randint(1, 100)                
                honeytoken_access_count = random.randint(0, total_access_count)
                shared_ip_count = random.randint(1, 10)

                avg_associated_malicious_score = round(random.random(), 3)  # 0.0 to 1.0

                # 2. Introduce a bias so that higher values => higher maliciousness
                # Normalize features into [0, 1] ranges
                norm_total_access = total_access_count / 100.0
                if total_access_count > 0:
                    norm_honeytoken = honeytoken_access_count / float(total_access_count)
                else:
                    norm_honeytoken = 0
                norm_shared_ip = shared_ip_count / 10.0

                # We'll combine these normalized features into a weighted sum
                # Weighted sum to emphasize some features more than others
                maliciousness = (
                    0.4 * norm_total_access +
                    0.3 * norm_honeytoken +
                    0.2 * norm_shared_ip +
                    0.1 * avg_associated_malicious_score
                )

                # 3. Add some noise so it's not perfectly correlated
                noise = random.uniform(-0.05, 0.05)
                maliciousness += noise

                # 4. Clamp value into [0, 1]
                maliciousness = max(0.0, min(1.0, maliciousness))

                actual_maliciousness_score = round(maliciousness, 3)

                # 5. Write the row
                data_writer.writerow([
                    total_access_count,
                    honeytoken_access_count,
                    shared_ip_count,
                    avg_associated_malicious_score,
                    actual_maliciousness_score
                ])

# Main function to generate synthetic data
if __name__ == '__main__':
    print('Generating synthetic data...')
    data_generator = GenerateData('data.csv')
    data_generator.generate_data()
    print('Data generation complete')