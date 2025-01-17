# Malicious User Detection System

## Overview
The Malicious User Detection System is a tool designed to track user interactions, identify fake or anomalous behaviors, and predict malicious users using AI models integrated with Neo4j. This project combines graph-based insights with machine learning to enhance system security and detect threats proactively.

## Features
- **Data Logging:** Tracks user actions, including endpoint accesses and honeytoken triggers.
- **Graph Database Integration:** Stores and analyzes user interaction data in Neo4j.
- **AI-Powered Predictions:** Leverages machine learning models to classify users as malicious or non-malicious based on interaction patterns.
- **Graph-Based Insights:** Detects complex relationships and behaviors using Neo4j's graph queries.
- **Visualization:** View relationships and malicious patterns in Neo4j Browser.

## How It Works
1. **Data Generation:** Synthetic data is created with labeled malicious and non-malicious users.
2. **AI Training:** A machine learning model is trained using labeled data to predict user behavior.
3. **Neo4j Storage:** User interactions are logged and stored in Neo4j as nodes and relationships.
4. **Insights Extraction:** Cypher queries analyze interaction patterns for AI input.
5. **Prediction:** The trained AI model predicts maliciousness, updating Neo4j with risk scores.
6. **Action:** The system blocks or flags high-risk users for further action.
