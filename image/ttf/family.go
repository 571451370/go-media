package ttf

type Family struct {
	Normal []byte
}

var VGA437 = Family{
	Normal: vga437,
}

var ProggyClean = Family{
	Normal: proggyClean,
}
