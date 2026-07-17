package server

// This file can be used by services to register their seeders via init() functions.
// Individual service packages should import db/seeders and call seeders.Register in their init.
// The blank imports here ensure those init() functions run when the server package is imported.
// Add new service seed registrations as needed.

import (
	// _ "abstractions/services/your_service" // example: triggers init registering seeders
	_ "abstractions/db/seeders" // ensure seeders package init runs
)
