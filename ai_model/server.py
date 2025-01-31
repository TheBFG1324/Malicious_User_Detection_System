# Flask server to load and train MLP model and predict maliciousness score

# Required Libraries
from flask import Flask, request, jsonify
from models.MLPRegressor import MaliciousnessPredictor

# Initialize Flask app
app = Flask(__name__)

# Load and train model
predictor = MaliciousnessPredictor('./models/data.csv')
predictor.train_model()

# Predict maliciousness score
@app.route('/predict', methods=['POST'])
def predict():
    # Get features from request
    features = request.get_json()
    print(features)
    total_access_count = features['total_access_count']
    honeytoken_access_count = features['honeytoken_access_count']
    shared_ip_count = features['shared_ip_count']
    avg_associated_malicious_score = features['avg_associated_malicious_score']

    # Predict maliciousness score
    predicted_maliciousness = predictor.predict_maliciousness([total_access_count, honeytoken_access_count, shared_ip_count, avg_associated_malicious_score])
    return jsonify({'maliciousness_score': predicted_maliciousness})

# Run Flask app
if __name__ == '__main__':
    app.run(debug=True)
