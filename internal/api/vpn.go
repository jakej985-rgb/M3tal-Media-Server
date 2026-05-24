package api

import (
	"encoding/json"
	"net/http"

	"github.com/jakej985-rgb/m3tal-core/internal/vpn"
)

// VPNHandlers provides HTTP handlers for Gluetun VPN Manager operations.
type VPNHandlers struct{}

// NewVPNHandlers creates a new instance of VPNHandlers.
func NewVPNHandlers() *VPNHandlers {
	return &VPNHandlers{}
}

// GetStatus returns the connection status and settings of the VPN.
// GET /api/v2/vpn/status
func (h *VPNHandlers) GetStatus(w http.ResponseWriter, r *http.Request) {
	mgr, err := vpn.NewManager()
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to initialize VPN manager: "+err.Error())
		return
	}

	status, err := mgr.GetStatus()
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	writeJSON(w, http.StatusOK, status)
}

// ControlVPN handles starting, stopping, or restarting the VPN container.
// POST /api/v2/vpn/control
func (h *VPNHandlers) ControlVPN(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Action string `json:"action"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid JSON body")
		return
	}

	mgr, err := vpn.NewManager()
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to initialize VPN manager: "+err.Error())
		return
	}

	var opErr error
	switch req.Action {
	case "start":
		opErr = mgr.Start()
	case "stop":
		opErr = mgr.Stop()
	case "restart":
		opErr = mgr.Restart()
	default:
		writeError(w, http.StatusBadRequest, "invalid action: must be 'start', 'stop', or 'restart'")
		return
	}

	if opErr != nil {
		writeError(w, http.StatusInternalServerError, opErr.Error())
		return
	}

	writeJSON(w, http.StatusOK, map[string]string{"status": "success", "action": req.Action})
}

// SwitchRegion handles switching the VPN region.
// POST /api/v2/vpn/region
func (h *VPNHandlers) SwitchRegion(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Region string `json:"region"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid JSON body")
		return
	}

	if req.Region == "" {
		writeError(w, http.StatusBadRequest, "region is required")
		return
	}

	mgr, err := vpn.NewManager()
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to initialize VPN manager: "+err.Error())
		return
	}

	if err := mgr.SwitchRegion(req.Region); err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	writeJSON(w, http.StatusOK, map[string]string{"status": "success", "region": req.Region})
}

// SyncPort handles manual port sync.
// POST /api/v2/vpn/sync-port
func (h *VPNHandlers) SyncPort(w http.ResponseWriter, r *http.Request) {
	mgr, err := vpn.NewManager()
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to initialize VPN manager: "+err.Error())
		return
	}

	port, err := mgr.SyncForwardedPort()
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	writeJSON(w, http.StatusOK, map[string]any{"status": "success", "forwarded_port": port})
}

// CheckLeak handles leak check.
// GET /api/v2/vpn/check-leak
func (h *VPNHandlers) CheckLeak(w http.ResponseWriter, r *http.Request) {
	mgr, err := vpn.NewManager()
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to initialize VPN manager: "+err.Error())
		return
	}

	isLeak, hostIP, vpnIP, err := mgr.CheckLeak()
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	// Automatic kill-switch trigger on leak detection
	var stoppedContainers []string
	if isLeak {
		stoppedContainers, _ = mgr.StopDependentContainers()
	}

	writeJSON(w, http.StatusOK, map[string]any{
		"leak":               isLeak,
		"host_ip":            hostIP,
		"vpn_ip":             vpnIP,
		"stopped_containers": stoppedContainers,
	})
}
