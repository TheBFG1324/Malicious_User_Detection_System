# This file trains and builds an NN model to predict the maliciousness score of a given user

# Required Libraries
from sklearn.neural_network import MLPRegressor
from sklearn.model_selection import train_test_split
from sklearn.preprocessing import MinMaxScaler
from sklearn.metrics import mean_squared_error, mean_absolute_error
import numpy as np

class MaliciousnessPredictor:
    def __init__(self, data_file):
        self.data_file = data_file
        self.model = None 
        self.scaler = None 

    def train_model(self):
        # Load data
        data = np.loadtxt(self.data_file, delimiter=',', skiprows=1)

        # Split data into features (X) and target (y)
        X = data[:, :-1] 
        y = data[:, -1]

        # Normalize features
        self.scaler = MinMaxScaler()
        X = self.scaler.fit_transform(X) 

        # Split into training and testing sets (80% train, 20% test)
        X_train, X_test, y_train, y_test = train_test_split(X, y, test_size=0.2, random_state=42)

        # Train neural network model
        self.model = MLPRegressor(hidden_layer_sizes=(100, 50), max_iter=1000, random_state=42)
        self.model.fit(X_train, y_train)

        # Evaluate model performance
        predictions = self.model.predict(X_test)

        mse = mean_squared_error(y_test, predictions)
        mae = mean_absolute_error(y_test, predictions)

        print(f"Model Evaluation:")
        print(f"Mean Squared Error (MSE): {mse:.4f}")
        print(f"Mean Absolute Error (MAE): {mae:.4f}")

    def predict_maliciousness(self, features):
        if self.model is None or self.scaler is None:
            raise ValueError("Model is not trained. Call train_model() first.")

        # Normalize input features using the same scaler as training
        features = self.scaler.transform([features])

        # Predict maliciousness score
        maliciousness_score = self.model.predict(features)[0]
        return maliciousness_score

if __name__ == "__main__":
    # Initialize the MaliciousnessPredictor with data file
    predictor = MaliciousnessPredictor("./data.csv")
    
    # Train the model
    predictor.train_model()

    # Example prediction
    features = [2,2,4,0.875]
    maliciousness_score = predictor.predict_maliciousness(features)
    print(f"Predicted Maliciousness Score: {maliciousness_score:.4f}")