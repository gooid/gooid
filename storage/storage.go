// Copyright 2018 The gooid Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.


package storage

import (
	"github.com/gooid/gooid/internal/ndk"
)

type AssetManager = app.AssetManager
type AssetDir = app.AssetDir
type Asset = app.Asset

/* Available modes for opening assets */
const (
	ASSET_MODE_UNKNOWN   = app.ASSET_MODE_UNKNOWN
	ASSET_MODE_RANDOM    = app.ASSET_MODE_RANDOM
	ASSET_MODE_STREAMING = app.ASSET_MODE_STREAMING
	ASSET_MODE_BUFFER    = app.ASSET_MODE_BUFFER
)
