// Copyright 2017 Google Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package media

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/godbus/dbus"
)

// name represents a dbus name that can be decomposed into an interface + member.
type name struct {
	iface  string
	member string
}

func (n name) String() string {
	return n.iface + "." + n.member
}

// buildMatchString builds a match string for the dbus (Add|Remove)Match methods.
func (n name) buildMatchString(sender string, args ...string) string {
	conditions := make([]string, 0)
	conditions = append(conditions, "type='signal'")
	conditions = append(conditions, fmt.Sprintf("interface='%s'", n.iface))
	conditions = append(conditions, fmt.Sprintf("member='%s'", n.member))
	if sender != "" {
		conditions = append(conditions, fmt.Sprintf("sender='%s'", sender))
	}
	for idx, val := range args {
		conditions = append(conditions, fmt.Sprintf("arg%d='%s'", idx, val))
	}
	return strings.Join(conditions, ",")
}

// Constants, signals and properties.
const (
	mprisInterface = "org.mpris.MediaPlayer2.Player"
	dbusInterface  = "org.freedesktop.DBus"
)

// Go doesn't support const structs.
var (
	// dbus methods
	methodNameHasOwner = name{dbusInterface, "NameHasOwner"}
	methodGetNameOwner = name{dbusInterface, "GetNameOwner"}
	methodAddMatch     = name{dbusInterface, "AddMatch"}
	methodRemoveMatch  = name{dbusInterface, "RemoveMatch"}

	// mpris methods
	mprisPlay      = name{mprisInterface, "Play"}
	mprisPause     = name{mprisInterface, "Pause"}
	mprisPlayPause = name{mprisInterface, "PlayPause"}
	mprisStop      = name{mprisInterface, "Stop"}
	mprisNext      = name{mprisInterface, "Next"}
	mprisPrev      = name{mprisInterface, "Previous"}
	mprisSeek      = name{mprisInterface, "Seek"}

	// mpris properties
	mprisRate     = name{mprisInterface, "Rate"}
	mprisPosition = name{mprisInterface, "Position"}
	mprisShuffle  = name{mprisInterface, "Shuffle"}
	mprisStatus   = name{mprisInterface, "PlaybackStatus"}
	mprisMetadata = name{mprisInterface, "Metadata"}

	// Dbus signals used for receiving updates about the media player.
	signalSeeked           = name{mprisInterface, "Seeked"}
	signalNameOwnerChanged = name{dbusInterface, "NameOwnerChanged"}
	signalPropChanged      = name{"org.freedesktop.DBus.Properties", "PropertiesChanged"}
)

// Some mpris players report numeric values as the wrong type. Fix that.
// TODO: See if this is a solved problem.

func getLong(l interface{}) int64 {
	switch l.(type) {
	case int, int8, int16, int32, int64:
		return reflect.ValueOf(l).Int()
	case uint, uint8, uint16, uint32, uint64:
		return int64(reflect.ValueOf(l).Uint())
	case float32, float64:
		return int64(reflect.ValueOf(l).Float())
	case dbus.Variant:
		return getLong(l.(dbus.Variant).Value())
	default:
		return 0
	}
}

func getDouble(d interface{}) float64 {
	switch d.(type) {
	case int, int8, int16, int32, int64:
		return float64(reflect.ValueOf(d).Int())
	case uint, uint8, uint16, uint32, uint64:
		return float64(reflect.ValueOf(d).Uint())
	case float32, float64:
		return reflect.ValueOf(d).Float()
	case dbus.Variant:
		return getDouble(d.(dbus.Variant).Value())
	default:
		return 0.0
	}
}
