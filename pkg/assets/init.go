// +build !js

package assets

func init() {
	p, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	DefaultAssetLoader = StaticFileAssetLoader{
		Pwd:  p,
		once: sync.Once{},
	}
}
