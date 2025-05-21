package neat

import (
	"math"
	"math/rand"
)

type Activation func(float64) float64

type NeuronGene struct {
	neuronId        int
	bias            float64
	activation      Activation
	hasBeencomputed bool
	value           float64
}

type LinkId struct {
	inputId  int
	outputId int
}

type LinkGene struct {
	linkId    LinkId
	weight    float64
	isEnabled bool
}

type Genome struct {
	genomeId   int
	numInputs  int
	numOutputs int
	neurons    []NeuronGene
	links      []LinkGene
}

func sigmoid(a float64) float64 {
	return 1.0 / (1.0 + math.Exp(-a))
}

func CreateGenome(genomeId, numInputs, numOutputs int) *Genome {
	return &Genome{genomeId: genomeId, numInputs: numInputs, numOutputs: numOutputs}
}

func (g *Genome) InitializeFromInitialConfig() {
	for i := 0; i < g.numInputs+g.numOutputs; i++ {
		newNeuron := NeuronGene{neuronId: i, activation: sigmoid}
		g.neurons = append(g.neurons, newNeuron)
	}

	for in := 0; in < g.numInputs; in++ {
		for out := g.numInputs; out < g.numInputs+g.numOutputs; out++ {
			newLink := LinkGene{
				linkId:    LinkId{inputId: in, outputId: out},
				weight:    (20 * rand.Float64()) - 10.0,
				isEnabled: true,
			}
			g.links = append(g.links, newLink)
		}
	}
}

func (g *Genome) LOLER() int {
	return len(g.neurons)
}

func (g *Genome) ComputeActivationLevels() {
	for k, v := range g.neurons {
		if k < g.numInputs {
		    v.hasBeencomputed = true
		} else {
		    v.hasBeencomputed = false
		}
		v.value = 0.0
	}

	for i := g.numInputs; i < g.numInputs + g.numOutputs; i++ {
		g.ComputeActivation(i)
	}
}

func (g *Genome) SetVisionInput(i int, d float64) {
	g.neurons[i].value = g.neurons[i].activation(d)
}

func (g *Genome) GetOutput(idx int) float64 {
	return g.neurons[g.numInputs+idx].value
}

func (g *Genome) Think() {
	g.ComputeActivationLevels()
}

func (g *Genome) ComputeActivation(neuronId int) (float64) {
	neuron := g.neurons[neuronId]

	if neuron.hasBeencomputed {
		return neuron.value
	}

	var res float64
	for _, link := range g.links {
		if link.isEnabled && link.linkId.outputId == neuron.neuronId {
			res += g.ComputeActivation(link.linkId.inputId)
		}
	}

	res = neuron.activation(res)
	neuron.hasBeencomputed = true
	neuron.value = res

	return res
}
