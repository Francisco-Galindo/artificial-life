package neat

import (
	"math"
	"math/rand"
	"slices"
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
	nextId     int
	numActiveNeurons int
	neurons    []*NeuronGene
	links      []*LinkGene
}

func sigmoid(a float64) float64 {
	return 1.0 / (1.0 + math.Exp(-a))
}

func relu(a float64) float64 {
	return math.Max(0, a)
}

func CreateGenome(genomeId, numInputs, numOutputs int) *Genome {
	return &Genome{genomeId: genomeId, numInputs: numInputs, numOutputs: numOutputs}
}

func (g *Genome) InitializeFromInitialConfig() {
	for i := 0; i < g.numInputs+g.numOutputs; i++ {
		newNeuron := NeuronGene{neuronId: i, activation: relu}
		g.neurons = append(g.neurons, &newNeuron)
	}

	for in := 0; in < g.numInputs; in++ {
		for out := g.numInputs; out < g.numInputs+g.numOutputs; out++ {
			newLink := LinkGene{
				linkId:    LinkId{inputId: in, outputId: out},
				weight:    max(-10, min(rand.NormFloat64()*2, 10)),
				isEnabled: true,
			}
			g.links = append(g.links, &newLink)
		}
	}
	g.numActiveNeurons = g.numInputs + g.numOutputs
}

func (g *Genome) ComputeActivationLevels() {
	queue := Queue{}

	for k := 0; k < len(g.neurons); k++ {
		v := g.neurons[k]
		if k < g.numInputs {
			queue.Enqueue(k)
			v.hasBeencomputed = true
		} else {
			v.hasBeencomputed = false
			v.value = 0.0
		}
	}

	for !queue.IsEmpty() {
		neuronId := queue.Dequeue()
		neuron := g.neurons[neuronId]

		if !neuron.hasBeencomputed {
			neuron.value = neuron.activation(neuron.value)
		}
		for _, link := range g.links {
			if link.isEnabled && link.linkId.inputId == neuron.neuronId {
				outIdx := link.linkId.outputId
				if !slices.Contains(queue.Elements, outIdx) {
				    queue.Enqueue(outIdx)
				}
				outNeuron := g.neurons[outIdx]
				outNeuron.value += neuron.value * link.weight
			}
		}
		neuron.hasBeencomputed = true
	}
}

func (g *Genome) SetVisionInput(i int, d float64) {
	g.neurons[i].value = g.neurons[i].activation(d)
}

func (g *Genome) SetHpInput(d int) {
	neuron := g.neurons[g.numInputs-1]
	g.neurons[g.numInputs-1].value = neuron.activation(float64(d))
}

func (g *Genome) GetOutput(idx int) float64 {
	return g.neurons[g.numInputs+idx].value
}

func (g *Genome) Think() {
	g.ComputeActivationLevels()
}

func (g *Genome) ComputeActivation(neuronId int) float64 {
	neuron := g.neurons[neuronId]

	if neuron.hasBeencomputed {
		return neuron.value
	}

	var res float64 = 0
	for _, link := range g.links {
		if link.isEnabled && link.linkId.outputId == neuron.neuronId {
			val := g.ComputeActivation(link.linkId.inputId)
			res += val * link.weight
		}
	}

	res = neuron.activation(res)
	neuron.hasBeencomputed = true
	neuron.value = res

	return res
}

func (g *Genome) Mutate() {
	g.mutateValues()
	g.mutateStructure()
}

func (g *Genome) MutateHighVariability() {
	for _, link := range g.links {
		if rand.Float64() < 0.30 {
			link.weight += max(-10, min(rand.NormFloat64()*0.25, 10))
		}
		if rand.Float64() < 0.01 {
			link.isEnabled = !link.isEnabled
		}
	}

	if rand.Float64() < 0.3 {
	    g.mutateAddLink()
	}
	if rand.Float64() < 0.2 {
		g.mutateRemoveLink()
	}
	if rand.Float64() < 0.1 {
		g.mutateAddNeuron()
	}
	if rand.Float64() < 0.02 {
		g.mutateRemoveNeuron()
	}
}

func (g *Genome) mutateStructure() {
	if rand.Float64() < 0.2 {
	    g.mutateAddLink()
	}
	if rand.Float64() < 0.1 {
		g.mutateRemoveLink()
	}
	if rand.Float64() < 0.05 {
		g.mutateAddNeuron()
	}
	if rand.Float64() < 0.01 {
		g.mutateRemoveNeuron()
	}
}

func (g *Genome) mutateRemoveNeuron() {
	n := len(g.neurons) - g.numOutputs - g.numInputs
	if n == 0 {
		return
	}

	hiddenId := rand.Intn(n) + g.numInputs + g.numOutputs
	for _, link := range g.links {
		if link.linkId.inputId == hiddenId || link.linkId.outputId == hiddenId {
			link.isEnabled = false
			link.weight = 1
		}
	}
	g.numActiveNeurons--
}

func (g *Genome) mutateAddNeuron() {
	if len(g.links) == 0 {
		return
	}

	k := rand.Intn(len(g.links))
	newId := len(g.neurons)
	newNeuron := NeuronGene{neuronId: newId, activation: relu}


	oldLink := g.links[k]
	oldLink.isEnabled = false

	g.links = append(g.links, &LinkGene{
		linkId: LinkId{inputId: oldLink.linkId.inputId, outputId: newId},
		weight: 1,
		isEnabled: true,
	})
	g.links = append(g.links, &LinkGene{
		linkId: LinkId{inputId: newId, outputId: oldLink.linkId.outputId},
		weight: oldLink.weight,
		isEnabled: true,
	})

	g.neurons = append(g.neurons, &newNeuron)
	g.numActiveNeurons++
}

func (g *Genome) mutateRemoveLink() {
	if len(g.links) == 0 {
		return
	}

	k := rand.Intn(len(g.links))
	g.links[k] = g.links[len(g.links)-1]
	g.links = g.links[:len(g.links)-1]
}

func (g *Genome) mutateAddLink() {
	inputId := rand.Intn(g.numInputs)
	outputId := rand.Intn(len(g.neurons) - g.numInputs) + g.numInputs
	for _, link := range g.links {
		linkId := link.linkId
		if linkId.inputId == inputId && linkId.outputId == outputId {
			link.isEnabled = true
			return
		}
	}

	g.links = append(g.links, &LinkGene{
		linkId: LinkId{inputId: inputId, outputId: outputId},
		weight: max(-10, min(rand.NormFloat64()*2, 10)),
		isEnabled: true,
	})

	visited := make([]bool, len(g.links))
	recStack := make([]bool, len(g.links))
	if g.hasCycle(inputId, visited, recStack) {
		g.links = g.links[:len(g.links)-1]
	}
}

func (g *Genome) hasCycle(cur int, visited, recStack []bool) bool {
	if recStack[cur] {
		return true
	}

	if visited[cur] {
		return false
	}

	visited[cur] = true
	recStack[cur] = true

	for _, link := range g.links {
		if link.linkId.inputId == cur {
			if g.hasCycle(link.linkId.outputId, visited, recStack) {
				return true
			}
		}
	}

	recStack[cur] = false
	return false
}

func (g *Genome) mutateValues() {
	for _, link := range g.links {
		if rand.Float64() < 0.10 {
			link.weight += max(-10, min(rand.NormFloat64()*0.25, 10))
		}
		if rand.Float64() < 0.01 {
			link.isEnabled = !link.isEnabled
		}
	}
}

func (g Genome) Copypapy() *Genome {
	return &g
}

func (g Genome) GetNumberOfNeurons() int {
	return g.numActiveNeurons
}
