package lighting

import (
	"image"
)

// Copied from golang/src/pkg/image/draw/draw.go
func clip(dst *image.RGBA, r *image.Rectangle, src *image.RGBA, sp *image.Point, mask *image.RGBA, mp *image.Point) {
	orig := r.Min
	*r = r.Intersect(dst.Bounds())
	*r = r.Intersect(src.Bounds().Add(orig.Sub(*sp)))
	if mask != nil {
		*r = r.Intersect(mask.Bounds().Add(orig.Sub(*mp)))
	}
	dx := r.Min.X - orig.X
	dy := r.Min.Y - orig.Y
	if dx == 0 && dy == 0 {
		return
	}
	(*sp).X += dx
	(*sp).Y += dy
	(*mp).X += dx
	(*mp).Y += dy
}

func DrawLightOverlay(dst *image.RGBA, r image.Rectangle, src *image.RGBA, sp image.Point, mask *image.RGBA, mp image.Point) {
	clip(dst, &r, src, &sp, mask, &mp)

	const m = 1<<16 - 1

	i0 := dst.PixOffset(r.Min.X, r.Min.Y)
	i1 := i0 + r.Dx()*4
	si0 := src.PixOffset(sp.X, sp.Y)
	mi0 := mask.PixOffset(mp.X, mp.Y)
	for y, sy, my := r.Min.Y, sp.Y, mp.Y; y != r.Max.Y; y, sy, my = y+1, sy+1, my+1 {
		for i, si, mi := i0, si0, mi0; i < i1; i, si, mi = i+4, si+4, mi+4 {
			ma := uint32(mask.Pix[mi+3])
			if ma == 0 {
				continue
			}
			ma |= ma << 8

			sr := uint32(src.Pix[si+0]) * 0x101
			sg := uint32(src.Pix[si+1]) * 0x101
			sb := uint32(src.Pix[si+2]) * 0x101
			sa := uint32(src.Pix[si+3]) * 0x101

			dr := uint32(dst.Pix[i+0])
			dg := uint32(dst.Pix[i+1])
			db := uint32(dst.Pix[i+2])
			da := uint32(dst.Pix[i+3])

			a := (m - (sa * ma / m)) * 0x101

			dst.Pix[i+0] = uint8((dr*a + sr*ma) / m >> 8)
			dst.Pix[i+1] = uint8((dg*a + sg*ma) / m >> 8)
			dst.Pix[i+2] = uint8((db*a + sb*ma) / m >> 8)
			dst.Pix[i+3] = uint8((da*a + sa*ma) / m >> 8)
		}
		i0 += dst.Stride
		i1 += dst.Stride
		si0 += src.Stride
		mi0 += mask.Stride
	}
}
