package api

// Sign is invoked by the User to notarize an artifact using the given functional options,
// By default, the artifact is notarized using status = meta.StatusTrusted, visibility meta.VisibilityPrivate.
func (u LcUser) Sign(artifact Artifact, options ...LcSignOption) (bool, uint64, error) {
	if artifact.Hash == "" {
		return false, 0, makeError("hash is missing", nil)
	}
	if artifact.Size < 0 {
		return false, 0, makeError("invalid size", nil)
	}

	o, err := makeLcSignOpts(options...)
	if err != nil {
		return false, 0, err
	}

	return u.createArtifact(artifact, o.status, o.attach)
}
