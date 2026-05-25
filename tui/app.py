import os
import sys
import time
from datetime import datetime

from textual.app import App, ComposeResult
from textual.widgets import Header, Footer, DataTable, RichLog, Static, TabbedContent, TabPane, Label, OptionList
from textual.containers import Container, Horizontal, Vertical
from textual.binding import Binding
from textual.reactive import reactive

# Add parent directory to sys.path to allow importing from tui sub-module
sys.path.insert(0, os.path.dirname(os.path.dirname(os.path.abspath(__file__))))

from tui.api_client import M3TALClient

class MetricsPanel(Static):
    """Bottom panel displaying real-time system metrics."""
    metrics_text = reactive("Loading system metrics...")

    def render(self) -> str:
        return self.metrics_text


class M3TALApp(App):
    TITLE = "M3TAL CONTROL CENTER (v2)"
    show_catalog = reactive(False)
    CSS = """
    Screen {
        background: #121212;
    }
    
    TabbedContent {
        background: #121212;
        margin: 0 1;
    }

    TabPane {
        padding: 0;
    }

    DataTable {
        background: #181818;
        color: #e0e0e0;
        border: rounded #2a2a2a;
        margin: 0 1;
        height: 1fr;
    }
    
    DataTable:focus {
        border: rounded #00E676;
    }

    OptionList {
        background: #181818;
        color: #e0e0e0;
        border: rounded #2a2a2a;
        margin: 0 1;
        height: 1fr;
    }

    OptionList:focus {
        border: rounded #00E676;
    }

    RichLog {
        background: #0d0d0d;
        color: #f0f0f0;
        border: rounded #2a2a2a;
        margin: 0 1;
        height: 1fr;
    }

    .pane-header {
        text-align: center;
        background: #1e1e1e;
        color: #00E676;
        text-style: bold;
        margin: 1 1 0 1;
        padding: 0 1;
        height: 1;
    }

    #metrics-panel {
        background: #1a1a1a;
        height: 3;
        border-top: hdouble #2a2a2a;
        padding: 0 2;
        content-align: left middle;
        color: #ffffff;
    }
    """

    BINDINGS = [
        Binding("q", "quit", "Quit"),
        Binding("r", "refresh", "Refresh Data"),
        Binding("u", "deploy_selected_stack", "Deploy Stack"),
        Binding("d", "stop_selected_stack", "Stop Stack"),
        Binding("s", "start_selected_container", "Start Container"),
        Binding("x", "stop_selected_container", "Stop Container"),
        Binding("t", "restart_selected_container", "Restart Container"),
        Binding("e", "toggle_selected_plugin", "Enable/Disable Plugin"),
        Binding("i", "install_selected_plugin", "Install Plugin"),
        Binding("c", "toggle_plugin_catalog", "Toggle Catalog/Installed"),
        Binding("k", "cancel_selected_job", "Cancel Job"),
    ]

    def __init__(self):
        super().__init__()
        self.client = M3TALClient()
        self.selected_stack = None
        self.selected_container = None
        self.log_tail_active = False
        self.last_log_check = {}

    def compose(self) -> ComposeResult:
        yield Header()
        
        with TabbedContent(initial="stacks-tab"):
            with TabPane("Stacks & Services", id="stacks-tab"):
                with Horizontal():
                    with Vertical():
                        yield Label("📁 STACKS", classes="pane-header")
                        yield DataTable(id="stacks-table")
                    with Vertical():
                        yield Label("🐳 STACK SERVICES", classes="pane-header")
                        yield DataTable(id="services-table")
            
            with TabPane("Container Logs", id="logs-tab"):
                with Horizontal():
                    with Vertical(id="logs-sidebar-container"):
                        yield Label("🔍 SELECT CONTAINER", classes="pane-header")
                        yield OptionList(id="logs-containers-list")
                    with Vertical():
                        yield Label("📄 LOG STREAM", classes="pane-header")
                        yield RichLog(id="log-viewer", highlight=True, markup=False)

            with TabPane("AI Queue & Models", id="ai-tab"):
                with Horizontal():
                    with Vertical():
                        yield Label("🤖 AI JOB QUEUE", classes="pane-header")
                        yield DataTable(id="ai-queue-table")
                    with Vertical():
                        yield Label("📦 OLLAMA MODELS", classes="pane-header")
                        yield OptionList(id="ai-models-list")

            with TabPane("Plugins Manager", id="plugins-tab"):
                yield Label("🧩 LOADED PLUGINS (ROUTES, STACKS, MIDDLEWARE)", id="plugins-header-label", classes="pane-header")
                yield DataTable(id="plugins-table")

        yield MetricsPanel(id="metrics-panel")
        yield Footer()

    def on_mount(self) -> None:
        # Initialize DataTables
        stacks_table = self.query_one("#stacks-table", DataTable)
        stacks_table.add_columns("Name", "Status", "Compose Path")
        stacks_table.cursor_type = "row"

        services_table = self.query_one("#services-table", DataTable)
        services_table.add_columns("Container Name", "Image", "State", "Ports")
        services_table.cursor_type = "row"

        ai_table = self.query_one("#ai-queue-table", DataTable)
        ai_table.add_columns("Job ID", "Type", "Priority", "Details", "Status")
        ai_table.cursor_type = "row"

        plugins_table = self.query_one("#plugins-table", DataTable)
        plugins_table.add_columns("Kind", "Name", "Version", "Author", "Status", "Warnings")
        plugins_table.cursor_type = "row"

        # Populate initial data and start periodic polling
        self.refresh_all()
        self.set_interval(2.0, self.refresh_all)

    def refresh_all(self) -> None:
        self.refresh_metrics()
        
        # Get active tab
        tabbed_content = self.query_one(TabbedContent)
        active_tab = tabbed_content.active
        
        if active_tab == "stacks-tab":
            self.refresh_stacks_and_services()
        elif active_tab == "logs-tab":
            self.refresh_logs_tab()
        elif active_tab == "ai-tab":
            self.refresh_ai_tab()
        elif active_tab == "plugins-tab":
            self.refresh_plugins_tab()

    def refresh_metrics(self) -> None:
        res = self.client.get_metrics()
        history_res = self.client.get_system_metrics_history()
        health_res = self.client.get_system_health()
        
        panel = self.query_one("#metrics-panel", MetricsPanel)
        
        if res.get("status") == "success" and "data" in res:
            metrics = res["data"]
            cpu = metrics.get("cpu_usage", 0.0)
            mem = metrics.get("memory_usage", 0.0)
            disk = metrics.get("disk_usage", 0.0)
            uptime_sec = metrics.get("uptime", 0)
            hostname = metrics.get("hostname", "unknown")
            
            # Format uptime
            days = uptime_sec // 86400
            hours = (uptime_sec % 86400) // 3600
            minutes = (uptime_sec % 3600) // 60
            uptime_str = f"{days}d {hours}h {minutes}m" if days > 0 else f"{hours}h {minutes}m"
            
            # Text progress bars
            def prog_bar(val):
                filled = int(val / 10)
                empty = 10 - filled
                bar = "█" * filled + "░" * empty
                return f"[{bar}] {val:.1f}%"

            # Generate sparklines
            history_list = []
            if history_res.get("status") == "success" and "data" in history_res:
                history_list = history_res["data"] or []
            
            def get_sparkline(key):
                if not history_list:
                    return ""
                chars = [' ', '▂', '▃', '▄', '▅', '▆', '▇', '█']
                # Keep last 15 ticks
                vals = [item.get(key, 0.0) for item in history_list][-15:]
                spark = []
                for v in vals:
                    idx = min(7, max(0, int(v / 12.5)))
                    spark.append(chars[idx])
                return " " + "".join(spark)
            
            # Health indicator
            health_str = "[bold #E0E0E0]?[/bold #E0E0E0]"
            if health_res.get("status") == "success" and "data" in health_res:
                health_data = health_res["data"]
                if health_data.get("status") == "healthy":
                    health_str = "[bold #00E676]🟢 OK[/bold #00E676]"
                else:
                    health_str = "[bold #FF1744]🔴 FAULT[/bold #FF1744]"

            panel.metrics_text = (
                f"💻 Host: [bold #00E676]{hostname}[/bold #00E676]  |  "
                f"Health: {health_str}  |  "
                f"Uptime: {uptime_str}  |  "
                f"CPU: {prog_bar(cpu)}{get_sparkline('cpu_usage')}  |  "
                f"RAM: {prog_bar(mem)}{get_sparkline('memory_usage')}  |  "
                f"Disk: {prog_bar(disk)}{get_sparkline('disk_usage')}"
            )
        else:
            panel.metrics_text = "[bold #FF1744]🔴 API Offline (Connection Error)[/bold #FF1744]"

    def refresh_stacks_and_services(self) -> None:
        # 1. Fetch stacks
        stacks_res = self.client.get_stacks()
        stacks_table = self.query_one("#stacks-table", DataTable)
        
        # Save current cursor row index
        current_stack_cursor = stacks_table.cursor_coordinate
        
        stacks_table.clear()
        stacks_list = []
        if stacks_res.get("status") == "success" and "data" in stacks_res:
            stacks_list = stacks_res["data"]
            for idx, stack in enumerate(stacks_list):
                name = stack.get("name", "unknown")
                status = stack.get("status", "discovered")
                path = stack.get("compose_path", "")
                
                # Format status color
                color = "#FFD600" # yellow
                if status == "running" or status == "success":
                    color = "#00E676" # green
                elif status == "failed" or status == "#FF1744":
                    color = "#FF1744" # red
                
                status_styled = f"[bold {color}]{status.upper()}[/bold {color}]"
                stacks_table.add_row(name, status_styled, path, key=name)
        
        # Restore cursor
        if current_stack_cursor and current_stack_cursor.row < len(stacks_list):
            stacks_table.cursor_coordinate = current_stack_cursor
        
        # Determine selected stack
        try:
            row_key, _ = stacks_table.coordinate_to_cell_key(stacks_table.cursor_coordinate)
            self.selected_stack = row_key.value
        except Exception:
            if stacks_list:
                self.selected_stack = stacks_list[0].get("name")
            else:
                self.selected_stack = None

        # 2. Fetch containers and filter by selected stack
        containers_res = self.client.get_containers()
        services_table = self.query_one("#services-table", DataTable)
        current_services_cursor = services_table.cursor_coordinate
        services_table.clear()
        
        if containers_res.get("status") == "success" and "data" in containers_res and self.selected_stack:
            containers = containers_res["data"]
            filtered_containers = []
            
            for c in containers:
                proj = c.get("labels", {}).get("com.docker.compose.project", "").lower()
                if proj == self.selected_stack.lower():
                    filtered_containers.append(c)
                    
            for c in filtered_containers:
                raw_name = c.get("names", ["unknown"])[0]
                name = raw_name.lstrip("/")
                image = c.get("image", "unknown")
                state = c.get("state", "unknown")
                status = c.get("status", "unknown")
                
                # Format Ports
                ports_list = c.get("ports", [])
                port_str = ", ".join([f"{p.get('public_port')}->{p.get('private_port')}" for p in ports_list if p.get("public_port")])
                
                state_color = "#FFD600"
                if state == "running":
                    state_color = "#00E676"
                elif state == "exited" or state == "dead":
                    state_color = "#FF1744"
                
                state_styled = f"[bold {state_color}]{state.upper()}[/bold {state_color}] ({status})"
                services_table.add_row(name, image, state_styled, port_str or "none", key=name)
            
            # Restore cursor
            if current_services_cursor and current_services_cursor.row < len(filtered_containers):
                services_table.cursor_coordinate = current_services_cursor
            
            try:
                row_key, _ = services_table.coordinate_to_cell_key(services_table.cursor_coordinate)
                self.selected_container = row_key.value
            except Exception:
                if filtered_containers:
                    self.selected_container = filtered_containers[0].get("names", [""])[0].lstrip("/")
                else:
                    self.selected_container = None

    def refresh_logs_tab(self) -> None:
        containers_res = self.client.get_containers()
        list_widget = self.query_one("#logs-containers-list", OptionList)
        
        current_index = list_widget.highlighted
        container_names = []
        
        if containers_res.get("status") == "success" and "data" in containers_res:
            containers = containers_res["data"]
            for c in containers:
                raw_name = c.get("names", ["unknown"])[0].lstrip("/")
                container_names.append(raw_name)
        
        # Populate sidebar option list
        list_widget.clear_options()
        for cname in container_names:
            list_widget.add_option(cname)
            
        if current_index is not None and current_index < len(container_names):
            list_widget.highlighted = current_index
            active_log_container = container_names[current_index]
        elif container_names:
            list_widget.highlighted = 0
            active_log_container = container_names[0]
        else:
            active_log_container = None
            
        if active_log_container:
            # Fetch and stream logs
            log_viewer = self.query_one("#log-viewer", RichLog)
            
            # Check if container changed
            if getattr(self, "_last_log_container", None) != active_log_container:
                log_viewer.clear()
                self._last_log_container = active_log_container
                self.last_log_check[active_log_container] = ""
                
            logs_res = self.client.get_container_logs(active_log_container, tail="100")
            if logs_res.get("status") == "success" and "data" in logs_res:
                logs_text = logs_res["data"]
                # Only write logs if they have changed or are new
                last_len = len(self.last_log_check.get(active_log_container, ""))
                if len(logs_text) > last_len:
                    new_logs = logs_text[last_len:]
                    self.last_log_check[active_log_container] = logs_text
                    for line in new_logs.splitlines():
                        log_viewer.write(line)
            else:
                log_viewer.write(f"Error loading logs: {logs_res.get('error', 'unknown error')}")

    def refresh_ai_tab(self) -> None:
        # 1. Fetch queue
        queue_res = self.client.get_queue()
        queue_table = self.query_one("#ai-queue-table", DataTable)
        current_cursor = queue_table.cursor_coordinate
        queue_table.clear()
        
        queue_list = []
        if queue_res.get("status") == "success" and "data" in queue_res:
            queue_list = queue_res["data"]
            for idx, job in enumerate(queue_list):
                jid = job.get("id", "unknown")
                jtype = job.get("type", "unknown")
                pVal = job.get("priority", 2)
                status = job.get("status", "pending")
                
                # Priority string
                prio_str = "Normal"
                if pVal == 3:
                    prio_str = "High"
                elif pVal == 1:
                    prio_str = "Low"

                # Extract details from payload
                payload = job.get("payload") or {}
                details = ""
                if jtype == "ai_generation":
                    details = payload.get("prompt", "")
                else:
                    details = str(payload)
                
                # Format status color
                color = "#FFD600" # yellow for pending/running
                if status == "completed":
                    color = "#00E676" # green
                elif status == "failed":
                    color = "#FF1744" # red
                    
                status_styled = f"[bold {color}]{status.upper()}[/bold {color}]"
                queue_table.add_row(jid, jtype, prio_str, details, status_styled, key=jid)
                
        if current_cursor and current_cursor.row < len(queue_list):
            queue_table.cursor_coordinate = current_cursor

        # 2. Fetch models list
        models_res = self.client.get_ai_models()
        models_list_widget = self.query_one("#ai-models-list", OptionList)
        models_list_widget.clear_options()
        
        if models_res.get("status") == "success" and "data" in models_res:
            for model in models_res["data"]:
                models_list_widget.add_option(model)

    def action_cancel_selected_job(self) -> None:
        queue_table = self.query_one("#ai-queue-table", DataTable)
        try:
            row_key, _ = queue_table.coordinate_to_cell_key(queue_table.cursor_coordinate)
            selected_jid = row_key.value
        except Exception:
            self.notify("No job selected in the queue table", severity="warning")
            return

        self.notify(f"✕ Cancelling Job: {selected_jid}...", severity="warning")
        res = self.client.cancel_queue_job(selected_jid)
        if res.get("status") == "success":
            self.notify(f"✅ Job {selected_jid} cancelled!", severity="success")
        else:
            self.notify(f"❌ Failed to cancel {selected_jid}: {res.get('error')}", severity="error")
        self.refresh_all()

    def refresh_plugins_tab(self) -> None:
        table = self.query_one("#plugins-table", DataTable)
        current_cursor = table.cursor_coordinate
        table.clear()
        
        # Update pane header label depending on whether we show catalog or installed
        try:
            label = self.query_one("#plugins-header-label", Label)
        except Exception:
            label = None
        
        if self.show_catalog:
            if label:
                label.update("🧩 PLUGINS CATALOG (AVAILABLE TO INSTALL / ENABLE)")
            res = self.client.get_plugin_catalog()
            rows_len = 0
            if res.get("status") == "success" and "data" in res:
                items = res["data"]
                rows_len = len(items)
                for item in items:
                    name = item.get("name", "unknown")
                    kind = item.get("kind", "Route")
                    version = item.get("version", "1.0.0")
                    author = item.get("author", "unknown")
                    status = item.get("status", "not_installed")
                    
                    if status == "enabled":
                        status_str = "[bold #00E676]ENABLED[/bold #00E676]"
                    elif status == "disabled":
                        status_str = "[bold #FFD600]DISABLED[/bold #FFD600]"
                    else:
                        status_str = "[bold #90a4ae]NOT INSTALLED[/bold #90a4ae]"
                        
                    warnings = ""
                    deps = [d.get("name") for d in item.get("dependencies", []) if d.get("required")]
                    if deps:
                        warnings = f"Requires: {', '.join(deps)}"
                    
                    table.add_row(kind, name, version, author, status_str, warnings, key=f"{kind}:{name}")
            
            # Restore cursor
            if current_cursor and current_cursor.row < rows_len:
                table.cursor_coordinate = current_cursor
        else:
            if label:
                label.update("🧩 LOADED PLUGINS (ROUTES, STACKS, MIDDLEWARE)")
            res = self.client.get_plugins()
            rows = []
            if res.get("status") == "success" and "data" in res:
                plugins = res["data"]
                for kind in ["routes", "stacks", "middleware"]:
                    items = plugins.get(kind, [])
                    for item in items:
                        meta = item.get("metadata", {})
                        name = meta.get("name", item.get("name", "unknown"))
                        version = meta.get("version", "1.0.0")
                        author = meta.get("author", "unknown")
                        enabled = item.get("enabled", False)
                        warnings_list = item.get("warnings", [])
                        
                        if enabled:
                            status_str = "[bold #00E676]ENABLED[/bold #00E676]"
                        else:
                            status_str = "[bold #FFD600]DISABLED[/bold #FFD600]"
                            
                        warnings_str = ""
                        if warnings_list:
                            warnings_str = f"[bold #FF1744]⚠️  {', '.join(warnings_list)}[/bold #FF1744]"
                        
                        kind_val = "Route"
                        if kind == "stacks":
                            kind_val = "Stack"
                        elif kind == "middleware":
                            kind_val = "Middleware"
                            
                        rows.append((kind_val, name, version, author, status_str, warnings_str))
                        
            for r in rows:
                table.add_row(*r, key=f"{r[0]}:{r[1]}")
                
            if current_cursor and current_cursor.row < len(rows):
                table.cursor_coordinate = current_cursor

    # --- Actions / Hotkeys ---
    def action_deploy_selected_stack(self) -> None:
        if self.selected_stack:
            self.notify(f"🚀 Deploying Stack: {self.selected_stack}...", severity="information")
            res = self.client.deploy_stack(self.selected_stack)
            if res.get("status") == "success":
                self.notify(f"✅ Stack {self.selected_stack} deployed successfully!", severity="success")
            else:
                self.notify(f"❌ Failed to deploy {self.selected_stack}: {res.get('error')}", severity="error")
            self.refresh_all()

    def action_stop_selected_stack(self) -> None:
        if self.selected_stack:
            self.notify(f"🧹 Stopping Stack: {self.selected_stack}...", severity="warning")
            res = self.client.stop_stack(self.selected_stack)
            if res.get("status") == "success":
                self.notify(f"✅ Stack {self.selected_stack} stopped successfully!", severity="success")
            else:
                self.notify(f"❌ Failed to stop {self.selected_stack}: {res.get('error')}", severity="error")
            self.refresh_all()

    def action_start_selected_container(self) -> None:
        if self.selected_container:
            self.notify(f"🐳 Starting container {self.selected_container}...", severity="information")
            res = self.client.control_container(self.selected_container, "start")
            if res.get("status") == "success":
                self.notify(f"✅ Container {self.selected_container} started!", severity="success")
            else:
                self.notify(f"❌ Failed to start {self.selected_container}: {res.get('error')}", severity="error")
            self.refresh_all()

    def action_stop_selected_container(self) -> None:
        if self.selected_container:
            self.notify(f"🛑 Stopping container {self.selected_container}...", severity="warning")
            res = self.client.control_container(self.selected_container, "stop")
            if res.get("status") == "success":
                self.notify(f"✅ Container {self.selected_container} stopped!", severity="success")
            else:
                self.notify(f"❌ Failed to stop {self.selected_container}: {res.get('error')}", severity="error")
            self.refresh_all()

    def action_restart_selected_container(self) -> None:
        if self.selected_container:
            self.notify(f"🔄 Restarting container {self.selected_container}...", severity="information")
            res = self.client.control_container(self.selected_container, "restart")
            if res.get("status") == "success":
                self.notify(f"✅ Container {self.selected_container} restarted!", severity="success")
            else:
                self.notify(f"❌ Failed to restart {self.selected_container}: {res.get('error')}", severity="error")
            self.refresh_all()

    def action_refresh(self) -> None:
        self.notify("Refreshing all dashboard views...", severity="information")
        self.refresh_all()

    def action_toggle_plugin_catalog(self) -> None:
        self.show_catalog = not self.show_catalog
        self.notify(f"Switched view to {'Catalog' if self.show_catalog else 'Installed Plugins'}", severity="information")
        self.refresh_all()

    def action_toggle_selected_plugin(self) -> None:
        table = self.query_one("#plugins-table", DataTable)
        try:
            row_key, _ = table.coordinate_to_cell_key(table.cursor_coordinate)
            selected_key = row_key.value
            kind, name = selected_key.split(":")
        except Exception:
            self.notify("No plugin selected", severity="warning")
            return

        row_index = table.cursor_coordinate.row
        status_cell = table.get_row_at(row_index)[4]
        
        is_enabled = "ENABLED" in status_cell
        is_disabled = "DISABLED" in status_cell
        is_not_installed = "NOT INSTALLED" in status_cell

        if is_not_installed:
            self.notify(f"Plugin {name} is not installed. Press 'i' to install it first.", severity="warning")
            return

        if is_enabled:
            self.notify(f"🛑 Disabling plugin {name}...", severity="warning")
            res = self.client.disable_plugin(name, kind)
            if res.get("status") == "success":
                self.notify(f"✅ Plugin {name} disabled!", severity="success")
            else:
                self.notify(f"❌ Failed to disable: {res.get('error')}", severity="error")
        else:
            self.notify(f"⚡ Enabling plugin {name}...", severity="information")
            res = self.client.enable_plugin(name, kind)
            if res.get("status") == "success":
                self.notify(f"✅ Plugin {name} enabled!", severity="success")
            else:
                err_msg = res.get('error', 'unknown error')
                if res.get("data") and "warnings" in res["data"]:
                    warnings = res["data"]["warnings"]
                    err_msg = f"Unsatisfied dependencies: {', '.join(warnings)}"
                self.notify(f"❌ Failed to enable: {err_msg}", severity="error")
                
        self.refresh_all()

    def action_install_selected_plugin(self) -> None:
        if not self.show_catalog:
            self.notify("Please switch to the Catalog view ('c') to install plugins", severity="warning")
            return
            
        table = self.query_one("#plugins-table", DataTable)
        try:
            row_key, _ = table.coordinate_to_cell_key(table.cursor_coordinate)
            selected_key = row_key.value
            kind, name = selected_key.split(":")
        except Exception:
            self.notify("No plugin selected", severity="warning")
            return

        row_index = table.cursor_coordinate.row
        status_cell = table.get_row_at(row_index)[4]
        
        if "NOT INSTALLED" not in status_cell:
            self.notify(f"Plugin {name} is already installed", severity="warning")
            return

        self.notify(f"📥 Installing plugin {name}...", severity="information")
        res = self.client.install_plugin(name, kind)
        if res.get("status") == "success":
            self.notify(f"✅ Plugin {name} installed successfully!", severity="success")
        else:
            self.notify(f"❌ Failed to install: {res.get('error')}", severity="error")
            
        self.refresh_all()

if __name__ == "__main__":
    app = M3TALApp()
    app.run()
