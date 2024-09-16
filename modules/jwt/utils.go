package jwt

// copy input map to output map
func mergeJwtClaims(input JwtClaim, output JwtClaim) {
	for k, v := range input {
		if _, ok := input[k]; ok {
			output[k] = v
		}
	}
}
