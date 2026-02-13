package game

// IGatekeeper defines the future service for cross-server travel.
type IGatekeeper interface {
	GetServerForSystem(systemID string) (string, error)
	RegisterUserTransfer(userID string, targetServerID string) (string, error)
}

// PlanetHazards defines environmental dangers (Gas, Heat, Gravity).
type PlanetHazards struct {
	Gravity       float64 `json:"gravity"`       // Default 1.0
	Atmosphere    string  `json:"atmosphere"`    // "Breathable", "Toxic", "Vacuum"
	ThermalRating float64 `json:"thermal_rating"`// -100 to +100
}

// WorldState represents dynamic faction influence.
type WorldState struct {
	FactionInfluence map[string]float64 `json:"faction_influence"`
}
