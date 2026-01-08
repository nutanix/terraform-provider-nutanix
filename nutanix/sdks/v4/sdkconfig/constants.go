package sdkconfig

// AllowVersionNegotiation controls whether Nutanix v4 generated SDK clients will attempt
// server-side API version negotiation.
//
// Keep this centralized so the behavior can be toggled across all v4 clients by changing
// a single value.
const AllowVersionNegotiation = true

// DefaultPort is the default Prism Central port used by v4 SDK clients when no port is provided.
const DefaultPort = 9440
