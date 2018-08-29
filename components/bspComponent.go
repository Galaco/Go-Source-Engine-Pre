package components

import (
	"github.com/galaco/go-me-engine/engine/interfaces"
	"github.com/go-gl/mathgl/mgl32"
	"github.com/galaco/go-me-engine/valve/bsp/tree"
	"github.com/galaco/go-me-engine/components/bsp"
	"github.com/galaco/go-me-engine/components/renderable"
	"github.com/galaco/bsp/primitives/visibility"
)

// BspComponent essentially extends a renderable component, as its large number
// of primitives require constant visibility culling; most effectively managed by
// this component itself
type BspComponent struct {
	RenderableComponent
	nodeTrees []tree.Node
	leafClusters map[int16][]*tree.Leaf

	cache []interfaces.IGPUMesh
	cachedPosition mgl32.Vec3
	currentLeaf *tree.Leaf
	visibilityLump *visibility.Vis
	faceList []interfaces.IPrimitive
}

// BspCompo
func (component *BspComponent) GetRenderables() []interfaces.IGPUMesh {
	return component.cache
}

func (component *BspComponent) UpdateVisibilityList(position mgl32.Vec3) {
	// View hasn't moved
	if position.ApproxEqual(component.cachedPosition) {
		return
	}
	// Still in the same node, visibility can only change when moving between nodes
	if component.currentLeaf != nil {
		if bsp.IsPointInLeaf(position, component.currentLeaf.Min, component.currentLeaf.Max) {
			return
		}
	}

	component.cachedPosition = position

	component.cache[0].(*renderable.GPUResourceDynamic).Reset()
	component.currentLeaf = bsp.FindCurrentLeaf(component.nodeTrees, component.cachedPosition)
	// If current == nil then we are outside the map. No visibility calculation
	if component.currentLeaf != nil && component.currentLeaf.ClusterId > -1 {
		faceList := bsp.BuildFaceListForVisibleClusters(
			component.nodeTrees,
			component.visibilityLump.GetVisibleIdsForCluster(component.currentLeaf.ClusterId))

		prims := make([]interfaces.IPrimitive, len(faceList))
		for idx,faceIdx := range faceList {
			prims[idx] = component.faceList[faceIdx]
		}

		component.cache[0].AddPrimitives(prims)

		// Shouldn't ever happen, but this is a catch all in a case
		// where fewer than 4 (smallest number of faces that can create a sealed volume
		// are apparently visible
		if len(prims) < 4 {
			component.cache[0].AddPrimitives(component.faceList)
		}
	} else {
		component.cache[0].AddPrimitives(component.faceList)
	}
}

func (component *BspComponent) recursiveBuildClusterList(node tree.INode) {
	if node.IsLeaf() {
		if component.leafClusters[node.(*tree.Leaf).ClusterId] == nil {
			component.leafClusters[node.(*tree.Leaf).ClusterId] = []*tree.Leaf{
				node.(*tree.Leaf),
			}
		} else {
			component.leafClusters[node.(*tree.Leaf).ClusterId] = append(
				component.leafClusters[node.(*tree.Leaf).ClusterId],
				node.(*tree.Leaf))
		}
	} else {
		for _,child := range node.(*tree.Node).Children {
			component.recursiveBuildClusterList(child)
		}
	}
}

func NewBspComponent(bspTrees []tree.Node, faceList []interfaces.IPrimitive, visibilityLump *visibility.Vis) *BspComponent{
	c := BspComponent{
		nodeTrees: bspTrees,
		leafClusters: map[int16][]*tree.Leaf{},
		visibilityLump: visibilityLump,
		cache: []interfaces.IGPUMesh{
			renderable.NewGPUResourceDynamic([]interfaces.IPrimitive{}),
		},
		faceList: faceList,
	}
	c.Etype = T_BspComponent

	for _, root := range c.nodeTrees {
		c.recursiveBuildClusterList(&root)
	}

	return &c
}