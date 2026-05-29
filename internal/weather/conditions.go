package weather

// descriptions maps WMO weather codes to short Spanish text.
var descriptions = map[int]string{
	0:  "Despejado",
	1:  "Mayormente despejado",
	2:  "Parcialmente nublado",
	3:  "Nublado",
	45: "Niebla",
	48: "Niebla con escarcha",
	51: "Llovizna ligera",
	53: "Llovizna moderada",
	55: "Llovizna densa",
	61: "Lluvia ligera",
	63: "Lluvia moderada",
	65: "Lluvia fuerte",
	71: "Nieve ligera",
	73: "Nieve moderada",
	75: "Nieve fuerte",
	80: "Chubascos ligeros",
	81: "Chubascos moderados",
	82: "Chubascos violentos",
	95: "Tormenta",
	96: "Tormenta con granizo ligero",
	99: "Tormenta con granizo fuerte",
}

// describe returns a Spanish description for a WMO weather code, or
// "Desconocido" when the code is not recognized.
func describe(code int) string {
	if d, ok := descriptions[code]; ok {
		return d
	}
	return "Desconocido"
}
