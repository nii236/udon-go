package assembly

type ArgTuple struct {
	UdonTypeName
	VarName
}

type EventArgs struct {
	Type UdonTypeName
	Name string
}

var EventTable = map[EventName][]*ArgTuple{
	"_start":              {},
	"_update":             {},
	"_lateUpdate":         {},
	"_interact":           {},
	"_fixedUpdate":        {},
	"_onAnimatorIk":       {{"Int32", "onAnimatorIkLayerIndex"}},
	"_onAnimatorMove":     {},
	"_onAudioFilterRead":  {{"SingleArray", "onAudioFilterReadData"}, {"onAudioFilterReadChannels", "Int32"}},
	"_onBecameInvisible":  {},
	"_onBecameVisible":    {},
	"_onCollisionEnter":   {{"Collision", "onCollisionEnterOther"}},
	"_onCollisionEnter2D": {{"Collision2D", "onCollisionEnter2DOther"}},
	"_onCollisionExit":    {{"Collision", "onCollisionExitOther"}},
	"_onCollisionExit2D":  {{"Collision2D", "onCollisionExit2DOther"}},
	"_onCollisionStay":    {{"Collision", "onCollisionStayOther"}},
	"_onCollisionStay2D":  {{"Collision2D", "onCollisionStay2DOther"}},
	// "_onControllerColliderHit": {{"ControllerColliderHit", "onControllerColliderHitHit"),// type ControllerColliderHit is not found}}
	"_onDestroy":                  {},
	"_onDisable":                  {},
	"_onDrawGizmos":               {},
	"_onDrawGizmosSelected":       {},
	"_onEnable":                   {},
	"_onGUI":                      {},
	"_onJointBreak":               {{"Single", "onJointBreakBreakForce"}},
	"_onJointBreak2D":             {{"Joint2D", "onJointBreak2DBrokenJoint"}},
	"_onMouseDown":                {},
	"_onMouseDrag":                {},
	"_onMouseEnter":               {},
	"_onMouseExit":                {},
	"_onMouseOver":                {},
	"_onMouseUp":                  {},
	"_onMouseUpAsButton":          {},
	"_onParticleCollision":        {{"GameObject", "onParticleCollisionOther"}},
	"_onParticleTrigger":          {},
	"_onPostRender":               {},
	"_onPreCull":                  {},
	"_onPreRender":                {},
	"_onRenderImage":              {{"RenderTexture", "onRenderImageSrc"}, {"RenderTexture", "onRenderImageDest"}},
	"_onRenderObject":             {},
	"_onTransformChildrenChanged": {},
	"_onTransformParentChanged":   {},
	"_onTriggerEnter":             {{"Collider", "onTriggerEnterOther"}},
	"_onTriggerEnter2D":           {{"Collider2D", "onTriggerEnter2DOther"}},
	"_onTriggerExit":              {{"Collider", "onTriggerExitOther"}},
	"_onTriggerExit2D":            {{"Collider2D", "onTriggerExit2DOther"}},
	"_onTriggerStay":              {{"Collider", "onTriggerStayOther"}},
	"_onTriggerStay2D":            {{"Collider2D", "onTriggerStay2DOther"}},
	"_onValidate":                 {},
	"_onWillRenderObject":         {},
	"_onDrop":                     {},
	"_onOwnershipTransferred":     {},
	"_onPickup":                   {},
	"_onPickupUseDown":            {},
	"_onPickupUseUp":              {},
	"_onPlayerJoined":             {{"VRCPlayerApi", "onPlayerJoinedPlayer"}},
	"_onPlayerLeft":               {{"VRCPlayerApi", "onPlayerLeftPlayer"}},
	"_onSpawn":                    {},
	"_onStationEntered":           {},
	"_onStationExited":            {},
	"_onVideoEnd":                 {},
	"_onVideoPause":               {},
	"_onVideoPlay":                {},
	"_onVideoStart":               {},
	"_onPreSerialization":         {},
	"_onDeserialization":          {},
	// Comment-outed event
	// onDataStorageAdded
	// onDataStorageChanged
	// onDataStorageRemoved
}
