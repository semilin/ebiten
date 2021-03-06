// Copyright 2019 The Ebiten Authors
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

package mipmap

import (
	"sync"

	"github.com/hajimehoshi/ebiten/internal/graphics"
)

var (
	theVerticesBackend = &verticesBackend{
		backend: make([]float32, graphics.VertexFloatNum*1024),
	}
)

type verticesBackend struct {
	backend []float32
	head    int
	m       sync.Mutex
}

func (v *verticesBackend) slice(n int, last bool) []float32 {
	v.m.Lock()

	need := n * graphics.VertexFloatNum
	if l := len(v.backend); v.head+need > l {
		for v.head+need > l {
			l *= 2
		}
		v.backend = make([]float32, l)
		v.head = 0
	}

	s := v.backend[v.head : v.head+need]
	if last {
		// If last is true, the vertices backend is sent to GPU and it is fine to reuse the slice.
		v.head = 0
	} else {
		v.head += need
	}

	v.m.Unlock()
	return s
}

func vertexSlice(n int, last bool) []float32 {
	return theVerticesBackend.slice(n, last)
}

func quadVertices(sx0, sy0, sx1, sy1 int, a, b, c, d, tx, ty float32, cr, cg, cb, ca float32, last bool) []float32 {
	x := float32(sx1 - sx0)
	y := float32(sy1 - sy0)
	ax, by, cx, dy := a*x, b*y, c*x, d*y
	u0, v0, u1, v1 := float32(sx0), float32(sy0), float32(sx1), float32(sy1)

	// This function is very performance-sensitive and implement in a very dumb way.
	vs := vertexSlice(4, last)
	_ = vs[:32]

	vs[0] = tx
	vs[1] = ty
	vs[2] = u0
	vs[3] = v0
	vs[4] = cr
	vs[5] = cg
	vs[6] = cb
	vs[7] = ca

	vs[8] = ax + tx
	vs[9] = cx + ty
	vs[10] = u1
	vs[11] = v0
	vs[12] = cr
	vs[13] = cg
	vs[14] = cb
	vs[15] = ca

	vs[16] = by + tx
	vs[17] = dy + ty
	vs[18] = u0
	vs[19] = v1
	vs[20] = cr
	vs[21] = cg
	vs[22] = cb
	vs[23] = ca

	vs[24] = ax + by + tx
	vs[25] = cx + dy + ty
	vs[26] = u1
	vs[27] = v1
	vs[28] = cr
	vs[29] = cg
	vs[30] = cb
	vs[31] = ca

	return vs
}
