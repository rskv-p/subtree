package subtree

import (
	"bytes"
)

//-------------------
// Function: genParts
//-------------------

// genParts breaks a filter subject (filter) into parts based on wildcards (`pwc '*'` or `fwc '>'`).
// It processes the input filter, identifies the wildcards, and separates the parts accordingly.
// Wildcards are used to separate the string into chunks, either as prefixes or suffixes.
// The function adds the parts into the `parts` slice and returns it.
func genParts(filter []byte, parts [][]byte) [][]byte {
	var start int
	for i, e := 0, len(filter)-1; i < len(filter); i++ {
		if filter[i] == tsep {
			// Case when the token is followed by a pwc (wildcard)
			if i < e && filter[i+1] == pwc && (i+2 <= e && filter[i+2] == tsep || i+1 == e) {
				if i > start {
					parts = append(parts, filter[start:i+1]) // Add part before pwc
				}
				parts = append(parts, filter[i+1:i+2]) // Add the pwc itself
				i++                                    // Skip pwc
				if i+2 <= e {
					i++ // Skip next tsep from the next part too.
				}
				start = i + 1
			} else if i < e && filter[i+1] == fwc && i+1 == e {
				// Case when we encounter an fwc (wildcard) at the end
				if i > start {
					parts = append(parts, filter[start:i+1]) // Add part before fwc
				}
				parts = append(parts, filter[i+1:i+2]) // Add the fwc itself
				i++                                    // Skip fwc
				start = i + 1
			}
		} else if filter[i] == pwc || filter[i] == fwc {
			// Wildcard must be at the start or preceded by tsep.
			if prev := i - 1; prev >= 0 && filter[prev] != tsep {
				continue
			}
			// Wildcard must be at the end or followed by tsep.
			if next := i + 1; next == e || next < e && filter[next] != tsep {
				continue
			}
			// We start with a pwc or fwc.
			parts = append(parts, filter[i:i+1])
			if i+1 <= e {
				i++ // Skip next tsep from next part too.
			}
			start = i + 1
		}
	}
	if start < len(filter) {
		// Check if we need to consume a leading tsep.
		if filter[start] == tsep {
			start++
		}
		parts = append(parts, filter[start:])
	}
	return parts
}

//-------------------
// Function: matchParts
//-------------------

// matchParts attempts to match the given parts against a fragment (frag), which could be a prefix for nodes or a suffix for leaves.
// It returns a modified list of parts, and a boolean indicating whether the match was successful or not.
func matchParts(parts [][]byte, frag []byte) ([][]byte, bool) {
	lf := len(frag)
	if lf == 0 {
		return parts, true // Empty fragment matches all parts
	}

	var si int
	lpi := len(parts) - 1

	for i, part := range parts {
		if si >= lf {
			return parts[i:], true // If we have consumed all of the fragment, return the remaining parts
		}
		lp := len(part)
		// Check for pwc or fwc placeholders.
		if lp == 1 {
			if part[0] == pwc {
				index := bytes.IndexByte(frag[si:], tsep)
				// If no tsep is found, it indicates we need to move to the next node from the caller.
				if index < 0 {
					if i == lpi {
						return nil, true
					}
					return parts[i:], true
				}
				si += index + 1
				continue
			} else if part[0] == fwc {
				// If we reach an fwc, we have matched the part.
				return nil, true
			}
		}
		end := min(si+lp, lf)
		// If part is bigger than the remaining fragment, adjust to a portion of the part.
		if si+lp > end {
			// Fragment is smaller than part itself.
			part = part[:end-si]
		}
		if !bytes.Equal(part, frag[si:end]) {
			return parts, false // If the part does not match the fragment, return false
		}
		// If there is still a portion of the fragment left, update and continue matching.
		if end < lf {
			si = end
			continue
		}
		// If we matched partially, do not move past the current part but update the part to what was consumed.
		if end < si+lp {
			if end >= lf {
				parts = append([][]byte{}, parts...) // Create a copy before modifying.
				parts[i] = parts[i][lf-si:]
			} else {
				i++
			}
			return parts[i:], true
		}
		if i == lpi {
			return nil, true
		}
		// If we have a wildcard gap, continue matching the next part up to the next tsep.
		si += len(part)
	}
	return parts, false
}
