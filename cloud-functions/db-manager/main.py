import functions_framework
from google.cloud import sqladmin_v1
from flask import jsonify
import os

PROJECT_ID = os.environ.get('GCP_PROJECT_ID', 'voicevox-custom')
INSTANCE_NAME = 'aso-db'

@functions_framework.http
def manage_db(request):
    """HTTP Cloud Function for managing Cloud SQL instance."""

    # Handle CORS preflight
    if request.method == 'OPTIONS':
        headers = {
            'Access-Control-Allow-Origin': '*',
            'Access-Control-Allow-Methods': 'GET, POST, OPTIONS',
            'Access-Control-Allow-Headers': 'Content-Type',
            'Access-Control-Max-Age': '3600'
        }
        return ('', 204, headers)

    headers = {
        'Access-Control-Allow-Origin': '*',
        'Content-Type': 'application/json'
    }

    action = request.args.get('action', 'status')

    try:
        client = sqladmin_v1.SqlInstancesServiceClient()
        instance_path = f"projects/{PROJECT_ID}/instances/{INSTANCE_NAME}"

        if action == 'start':
            return start_instance(client, headers)
        elif action == 'stop':
            return stop_instance(client, headers)
        elif action == 'status':
            return get_status(client, headers)
        else:
            return (jsonify({'error': 'Invalid action. Use: start, stop, or status'}), 400, headers)

    except Exception as e:
        return (jsonify({'error': str(e)}), 500, headers)


def start_instance(client, headers):
    """Start the Cloud SQL instance."""
    request = sqladmin_v1.SqlInstancesGetRequest(
        project=PROJECT_ID,
        instance=INSTANCE_NAME
    )
    instance = client.get(request=request)

    if instance.settings.activation_policy == sqladmin_v1.Settings.SqlActivationPolicy.ALWAYS:
        return (jsonify({
            'status': 'already_running',
            'message': 'Database is already running'
        }), 200, headers)

    # Update activation policy to ALWAYS (start)
    instance.settings.activation_policy = sqladmin_v1.Settings.SqlActivationPolicy.ALWAYS

    update_request = sqladmin_v1.SqlInstancesUpdateRequest(
        project=PROJECT_ID,
        instance=INSTANCE_NAME,
        body=instance
    )

    client.update(request=update_request)

    return (jsonify({
        'status': 'starting',
        'message': 'Database start initiated'
    }), 200, headers)


def stop_instance(client, headers):
    """Stop the Cloud SQL instance."""
    request = sqladmin_v1.SqlInstancesGetRequest(
        project=PROJECT_ID,
        instance=INSTANCE_NAME
    )
    instance = client.get(request=request)

    if instance.settings.activation_policy == sqladmin_v1.Settings.SqlActivationPolicy.NEVER:
        return (jsonify({
            'status': 'already_stopped',
            'message': 'Database is already stopped'
        }), 200, headers)

    # Update activation policy to NEVER (stop)
    instance.settings.activation_policy = sqladmin_v1.Settings.SqlActivationPolicy.NEVER

    update_request = sqladmin_v1.SqlInstancesUpdateRequest(
        project=PROJECT_ID,
        instance=INSTANCE_NAME,
        body=instance
    )

    client.update(request=update_request)

    return (jsonify({
        'status': 'stopping',
        'message': 'Database stop initiated'
    }), 200, headers)


def get_status(client, headers):
    """Get the current status of the Cloud SQL instance."""
    request = sqladmin_v1.SqlInstancesGetRequest(
        project=PROJECT_ID,
        instance=INSTANCE_NAME
    )
    instance = client.get(request=request)

    is_running = instance.settings.activation_policy == sqladmin_v1.Settings.SqlActivationPolicy.ALWAYS

    return (jsonify({
        'status': 'running' if is_running else 'stopped',
        'state': str(instance.state),
        'activation_policy': str(instance.settings.activation_policy)
    }), 200, headers)
