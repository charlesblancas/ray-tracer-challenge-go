package main

import (
	"math"
	"reflect"
	"testing"
)

func TestCreateWorld(t *testing.T) {
	w := World{}

	if len(w.Lights) != 0 {
		t.Errorf("Expected there to be no light source but there is %v", w.Lights[0])
	}

	if len(w.Objects) != 0 {
		t.Errorf("Expected there to be no objects but there are %v", w.Objects)
	}
}

func TestDefaultWorld(t *testing.T) {
	l, err := pointLight(point(-10, 10, -10), Color{1, 1, 1})
	if err != nil {
		t.Fatal(err)
	}
	s1 := sphere()
	m := material()
	m.Color = Color{0.8, 1.0, 0.6}
	m.Diffuse = 0.7
	m.Specular = 0.2
	s1.Material = m

	s2 := sphere()
	s2.Transform = scaling(0.5, 0.5, 0.5)

	w, err := defaultWorld()
	if err != nil {
		t.Fatal(err)
	}

	expected := World{[]Sphere{s1, s2}, []PointLight{l}}

	if !reflect.DeepEqual(w.Lights[0], l) ||
		!reflect.DeepEqual(w.Objects[0], s1) ||
		!reflect.DeepEqual(w.Objects[1], s2) {
		t.Errorf("Error with default world %v is supposed to be %v", w, expected)
	}
}

func TestIntersectWorldWithRay(t *testing.T) {
	w, err := defaultWorld()
	if err != nil {
		t.Fatal(err)
	}
	r, err := ray(point(0, 0, -5), vector(0, 0, 1))
	if err != nil {
		t.Fatal(err)
	}

	xs, err := worldRayIntersect(w, r)
	if err != nil {
		t.Fatal(err)
	}
	if len(xs) != 4 ||
		!floatEqual(xs[0].t, 4) ||
		!floatEqual(xs[1].t, 4.5) ||
		!floatEqual(xs[2].t, 5.5) ||
		!floatEqual(xs[3].t, 6) {
		t.Errorf("Expected ts to be [4,4.5,5.5,6] but intersection array is %v", xs)
	}
}

func TestPrecomputeStateOfIntersection(t *testing.T) {
	r, err := ray(point(0, 0, -5), vector(0, 0, 1))
	if err != nil {
		t.Fatal(err)
	}

	shape := sphere()
	i := Intersection{shape, 4}
	comps, err := prepareComputations(i, r)
	if err != nil {
		t.Fatal(err)
	}

	expected := Computation{i.Object, i.t, point(0, 0, -1), vector(0, 0, -1), vector(0, 0, -1), false}

	if !floatEqual(comps.t, expected.t) ||
		!reflect.DeepEqual(comps.Object, expected.Object) ||
		!tupleEqual(comps.Point, expected.Point) ||
		!tupleEqual(comps.EyeV, expected.EyeV) ||
		!tupleEqual(comps.NormalV, expected.NormalV) {
		t.Errorf("Expected %v to equal %v", comps, expected)
	}
}

func TestHitWhenIntersectionIsOutside(t *testing.T) {
	r, err := ray(point(0, 0, -5), vector(0, 0, 1))
	if err != nil {
		t.Fatal(err)
	}
	shape := sphere()
	i := Intersection{shape, 4}
	comps, err := prepareComputations(i, r)
	if err != nil {
		t.Fatal(err)
	}
	if comps.IsInside {
		t.Errorf("Expected IsInside to be false but it is true")
	}
}

func TestHitWhenIntersectionIsInside(t *testing.T) {
	r, err := ray(point(0, 0, 0), vector(0, 0, 1))
	if err != nil {
		t.Fatal(err)
	}
	shape := sphere()
	i := Intersection{shape, 1}
	comps, err := prepareComputations(i, r)
	if err != nil {
		t.Fatal(err)
	}

	expected := Computation{i.Object, i.t, point(0, 0, 1), vector(0, 0, -1), vector(0, 0, -1), true}

	if !floatEqual(comps.t, expected.t) ||
		!tupleEqual(comps.Point, expected.Point) ||
		!tupleEqual(comps.EyeV, expected.EyeV) ||
		!tupleEqual(comps.NormalV, expected.NormalV) ||
		!comps.IsInside {
		t.Errorf("Expected %v to equal %v", comps, expected)
	}
}

func TestShadeAnIntersection(t *testing.T) {
	w, err := defaultWorld()
	if err != nil {
		t.Fatal(err)
	}
	r, err := ray(point(0, 0, -5), vector(0, 0, 1))
	if err != nil {
		t.Fatal(err)
	}
	shape := w.Objects[0]
	i := Intersection{shape, 4}
	comps, err := prepareComputations(i, r)
	if err != nil {
		t.Fatal(err)
	}
	c := shadeHit(w, comps)
	expected := Color{0.38066, 0.47583, 0.2855}

	if !colorEqual(c, expected) {
		t.Errorf("%v not equal to %v", c, expected)
	}
}

func TestShadeAnIntersectionFromInside(t *testing.T) {
	w, err := defaultWorld()
	if err != nil {
		t.Fatal(err)
	}
	l, err := pointLight(point(0, 0.25, 0), Color{1, 1, 1})
	if err != nil {
		t.Fatal(err)
	}
	w.Lights[0] = l
	r, err := ray(point(0, 0, 0), vector(0, 0, 1))
	if err != nil {
		t.Fatal(err)
	}
	shape := w.Objects[1]
	i := Intersection{shape, 0.5}
	comps, err := prepareComputations(i, r)
	if err != nil {
		t.Fatal(err)
	}
	c := shadeHit(w, comps)
	expected := Color{0.90498, 0.90498, 0.90498}

	if !colorEqual(c, expected) {
		t.Errorf("%v not equal to %v", c, expected)
	}
}

func TestColorWhenRayMisses(t *testing.T) {
	w, err := defaultWorld()
	if err != nil {
		t.Fatal(err)
	}
	r, err := ray(point(0, 0, -5), vector(0, 1, 0))
	if err != nil {
		t.Fatal(err)
	}
	c, err := colorAt(w, r)
	if err != nil {
		t.Fatal(err)
	}
	expected := Color{0, 0, 0}
	if !colorEqual(c, expected) {
		t.Errorf("Expected %v to equal %v", c, expected)
	}
}

func TestColorWhenRayHits(t *testing.T) {
	w, err := defaultWorld()
	if err != nil {
		t.Fatal(err)
	}
	r, err := ray(point(0, 0, -5), vector(0, 0, 1))
	if err != nil {
		t.Fatal(err)
	}
	c, err := colorAt(w, r)
	if err != nil {
		t.Fatal(err)
	}
	expected := Color{0.38066, 0.47583, 0.2855}
	if !colorEqual(c, expected) {
		t.Errorf("Expected %v to equal %v", c, expected)
	}
}

func TestColorWhenIntersectionBehindRay(t *testing.T) {
	w, err := defaultWorld()
	if err != nil {
		t.Fatal(err)
	}
	w.Objects[0].Material.Ambient = 1.0
	w.Objects[1].Material.Ambient = 1.0
	r, err := ray(point(0, 0, 0.75), vector(0, 0, -1))
	if err != nil {
		t.Fatal(err)
	}
	c, err := colorAt(w, r)
	if err != nil {
		t.Fatal(err)
	}
	expected := w.Objects[1].Material.Color
	if !colorEqual(c, expected) {
		t.Errorf("Expected %v to equal %v", c, expected)
	}
}

func TestTransformationMatrixForDefaultOrientation(t *testing.T) {
	tr, err := viewTransform(point(0, 0, 0), point(0, 0, -1), vector(0, 1, 0))
	if err != nil {
		t.Fatal(err)
	}
	expected := matrixConstructIdentity(4)
	if !matrixEqual(tr, expected) {
		t.Errorf("Expected %v to equal %v", tr, expected)
	}
}

func TestTransformationMatrixInPositiveZ(t *testing.T) {
	tr, err := viewTransform(point(0, 0, 0), point(0, 0, 1), vector(0, 1, 0))
	if err != nil {
		t.Fatal(err)
	}
	expected := scaling(-1, 1, -1)
	if !matrixEqual(tr, expected) {
		t.Errorf("Expected %v to equal %v", tr, expected)
	}
}

func TestViewTransformationMovesTheWorld(t *testing.T) {
	tr, err := viewTransform(point(0, 0, 8), point(0, 0, 0), vector(0, 1, 0))
	if err != nil {
		t.Fatal(err)
	}
	expected := translation(0, 0, -8)
	if !matrixEqual(tr, expected) {
		t.Errorf("Expected %v to equal %v", tr, expected)
	}
}

func TestArbitraryViewTransformation(t *testing.T) {
	tr, err := viewTransform(point(1, 3, 2), point(4, -2, 8), vector(1, 1, 0))
	if err != nil {
		t.Fatal(err)
	}
	expected := matrixConstruct([][]float64{
		{-0.50709, 0.50709, 0.67612, -2.36643},
		{0.76772, 0.60609, 0.12122, -2.82843},
		{-0.35857, 0.59761, -0.71714, 0.00000},
		{0.00000, 0.00000, 0.00000, 1.00000},
	})
	if !matrixEqual(tr, expected) {
		t.Errorf("Expected %v to equal %v", tr, expected)
	}
}

func TestConstructCamera(t *testing.T) {
	cam := camera(160, 120, math.Pi/2.)
	expected := Camera{matrixConstructIdentity(4), 160, 120, math.Pi / 2., -1, -1, -1}
	if cam.HSize != expected.HSize ||
		cam.VSize != expected.VSize ||
		!floatEqual(cam.FieldOfView, expected.FieldOfView) ||
		!matrixEqual(cam.Transform, expected.Transform) {
		t.Errorf("Expected %v to equal %v", cam, expected)
	}
}

func TestPixelSizeForHorizontalCanvas(t *testing.T) {
	c := camera(200, 125, math.Pi/2.)
	expected := 0.01
	if !floatEqual(c.PixelSize, expected) {
		t.Errorf("Expected %v to equal %v", c.PixelSize, expected)
	}
}

func TestPixelSizeForVerticalCanvas(t *testing.T) {
	c := camera(125, 200, math.Pi/2.)
	expected := 0.01
	if !floatEqual(c.PixelSize, expected) {
		t.Errorf("Expected %v to equal %v", c.PixelSize, expected)
	}
}

func TestConstructRayThroughCenterOfTheCanvas(t *testing.T) {
	c := camera(201, 101, math.Pi/2.)
	r, err := rayForPixel(c, 100, 50)
	if err != nil {
		t.Fatal(err)
	}
	e1 := point(0, 0, 0)
	if !tupleEqual(r.Origin, e1) {
		t.Errorf("Expected %v to equal %v", r.Origin, e1)
	}
	e2 := vector(0, 0, -1)
	if !tupleEqual(r.Direction, e2) {
		t.Errorf("Expected %v to equal %v", r.Origin, e2)
	}
}

func TestConstructRayThroughCornerOfTheCanvas(t *testing.T) {
	c := camera(201, 101, math.Pi/2.)
	r, err := rayForPixel(c, 0, 0)
	if err != nil {
		t.Fatal(err)
	}
	e1 := point(0, 0, 0)
	if !tupleEqual(r.Origin, e1) {
		t.Errorf("Expected %v to equal %v", r.Origin, e1)
	}
	e2 := vector(0.66519, 0.33259, -0.66851)
	if !tupleEqual(r.Direction, e2) {
		t.Errorf("Expected %v to equal %v", r.Origin, e2)
	}
}

func TestConstructRayWhenCameraIsTransformed(t *testing.T) {
	c := camera(201, 101, math.Pi/2.)
	tr, err := matrix4x4Multiply(rotationY(math.Pi/4), translation(0, -2, 5))
	if err != nil {
		t.Fatal(err)
	}
	c.Transform = tr
	r, err := rayForPixel(c, 100, 50)
	if err != nil {
		t.Fatal(err)
	}
	e1 := point(0, 2, -5)
	if !tupleEqual(r.Origin, e1) {
		t.Errorf("Expected %v to equal %v", r.Origin, e1)
	}
	e2 := vector(math.Sqrt2/2., 0, -math.Sqrt2/2.)
	if !tupleEqual(r.Direction, e2) {
		t.Errorf("Expected %v to equal %v", r.Origin, e2)
	}
}

func TestRenderWorldWithCamera(t *testing.T) {
	w, err := defaultWorld()
	if err != nil {
		t.Fatal(err)
	}
	c := camera(11, 11, math.Pi/2.)
	from := point(0, 0, -5)
	to := point(0, 0, 0)
	up := vector(0, 1, 0)
	tr, err := viewTransform(from, to, up)
	if err != nil {
		t.Fatal(err)
	}
	c.Transform = tr
	image, err := render(c, w)
	if err != nil {
		t.Fatal(err)
	}
	wanted := pixelAt(image, 5, 5)
	expected := Color{0.38066, 0.47583, 0.2855}
	if !colorEqual(wanted, expected) {
		t.Errorf("Expected %v to equal %v", wanted, expected)
	}
}
