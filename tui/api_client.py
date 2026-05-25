import os
import requests
from urllib.parse import urljoin

class M3TALClient:
    def __init__(self):
        self.api_url = os.environ.get("API_URL", "http://localhost:5050")
        self.api_token = os.environ.get("API_TOKEN", "m3tal-secret-token")

    def _get_headers(self):
        headers = {}
        if self.api_token:
            headers["X-API-Token"] = self.api_token
        return headers

    def _request(self, method, path, json_data=None, params=None):
        url = urljoin(self.api_url.rstrip("/") + "/", path.lstrip("/"))
        try:
            resp = requests.request(
                method, 
                url, 
                headers=self._get_headers(), 
                json=json_data, 
                params=params,
                timeout=5
            )
            if resp.status_code == 401:
                return {"status": "error", "error": "Unauthorized API Access"}
            
            data = resp.json()
            # Standard M3TAL API responses are wrapped: {"status": "success/error", "data": ...}
            if isinstance(data, dict) and "status" in data:
                return data
            return {"status": "success", "data": data}
        except Exception as e:
            return {"status": "error", "error": str(e)}

    def get_stacks(self):
        return self._request("GET", "/api/v2/stacks")

    def deploy_stack(self, name):
        return self._request("POST", f"/api/v2/stacks/{name}/up")

    def stop_stack(self, name):
        return self._request("POST", f"/api/v2/stacks/{name}/down")

    def get_containers(self):
        # /api/containers returns all containers
        return self._request("GET", "/api/containers")

    def control_container(self, name, action):
        # action is 'start', 'stop', 'restart'
        return self._request("POST", f"/api/containers/{action}", json_data={"name": name})

    def get_container_logs(self, name, tail="100"):
        return self._request("GET", f"/api/v2/containers/{name}/logs", params={"tail": tail})

    def get_metrics(self):
        # /api/metrics returns system metrics (or /api/v2/system/metrics)
        return self._request("GET", "/api/metrics")

    def get_ai_queue(self):
        return self._request("GET", "/api/v2/ai/queue")

    def get_ai_models(self):
        return self._request("GET", "/api/v2/ai/models")

    def get_plugins(self):
        return self._request("GET", "/api/v2/plugins")

    def enable_plugin(self, name, kind):
        return self._request("POST", "/api/v2/plugins/enable", json_data={"name": name, "kind": kind})

    def disable_plugin(self, name, kind):
        return self._request("POST", "/api/v2/plugins/disable", json_data={"name": name, "kind": kind})

    def install_plugin(self, name, kind):
        return self._request("POST", "/api/v2/plugins/install", json_data={"name": name, "kind": kind})

    def uninstall_plugin(self, name, kind):
        return self._request("POST", "/api/v2/plugins/uninstall", json_data={"name": name, "kind": kind})

    def get_plugin_catalog(self):
        return self._request("GET", "/api/v2/plugins/catalog")
