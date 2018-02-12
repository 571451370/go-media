package ttf

type Family struct {
	Normal []byte
	Light  []byte
	Bold   []byte
}

var VGA437 = Family{
	Normal: vga437,
}

var ProggyClean = Family{
	Normal: proggyClean,
}

var Roboto = Family{
	Normal: robotoNormal,
	Light:  robotoLight,
	Bold:   robotoBold,
}
