# This file trains and builds an NN model to predict the maliciousness score of a given user

# Required Libraries
from sklearn.neural_network import MLPRegressor
from sklearn.model_selection import train_test_split
from sklearn.preprocessing import MinMaxScaler
import numpy as np

class MaliciousnessPredictor:
    def __init__(self, data_file):
        self.data_file = data_file
        self.model = None 
        self.scaler = None 
        self.target_scaler = None  # New scaler for y

    def train_model(self):
        # Load data
        data = np.loadtxt(self.data_file, delimiter=',', skiprows=1)

        # Split data into features (X) and target (y)
        X = data[:, :-1] 
        y = data[:, -1].reshape(-1, 1)  # Reshape for MinMaxScaler

        # Normalize features
        self.scaler = MinMaxScaler()
        X = self.scaler.fit_transform(X)

        # Normalize target (y)
        self.target_scaler = MinMaxScaler()
        y = self.target_scaler.fit_transform(y).ravel()  # Flatten y after scaling

        # Split into training and testing sets (80% train, 20% test)
        X_train, X_test, y_train, y_test = train_test_split(X, y, test_size=0.2, random_state=42)

        # Train neural network model
        self.model = MLPRegressor(hidden_layer_sizes=(100, 50), max_iter=1000, random_state=42)
        self.model.fit(X_train, y_train)

    def predict_maliciousness(self, features):
        if self.model is None or self.scaler is None or self.target_scaler is None:
            raise ValueError("Model is not trained. Call train_model() first.")

        # Normalize input features using the same scaler as training
        features = self.scaler.transform([features])

        # Predict maliciousness score
        scaled_prediction = self.model.predict(features)[0]

        # Inverse transform prediction to original range
        maliciousness_score = self.target_scaler.inverse_transform([[scaled_prediction]])[0, 0]
        maliciousness_score = np.clip(maliciousness_score, 0, 1)
        
        return maliciousness_score
