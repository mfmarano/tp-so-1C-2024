package responses

type IniciarProcesoResponse struct {
	Pid int `json:"pid"`
}

type EstadoProcesoResponse struct {
	State string `json:"state"`
}

type ProcesoResponse struct {
	Pid   int    `json:"pid"`
	State string `json:"state"`
}
