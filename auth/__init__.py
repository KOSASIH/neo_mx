# auth/__init__.py
import os
import json
from flask import Flask, request, jsonify
from authlib.integrations.flask_client import OAuth

app = Flask(__name__)
app.config['SECRET_KEY'] = os.urandom(24)

oauth = OAuth(app)

oauth.register(
    name='neo_mx',
    client_id='your_client_id',
    client_secret='your_client_secret',
    access_token_url='https://your_oidc_provider.com/token',
    access_token_params=None,
    authorize_url='https://your_oidc_provider.com/authorize',
    authorize_params=None,
    api_base_url='https://your_oidc_provider.com/',
    client_kwargs={'scope': 'openid profile email'}
)

@app.route('/login')
def login():
    return oauth.neo_mx.authorize_redirect(request)

@app.route('/logout')
def logout():
    oauth.neo_mx.logout()
    return 'Logged out'

@app.route('/api/data')
def get_data():
    token = oauth.neo_mx.authorize_access_token()
    if token:
        # Use the token to fetch data from the API
        return jsonify({'data': [...]})
    return 'Unauthorized', 401
