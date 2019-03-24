package propertyimage

import (
	"fmt"
	vocab "github.com/go-fed/activity/streams/vocab"
	"net/url"
)

// ActivityStreamsImagePropertyIterator is an iterator for a property. It is
// permitted to be one of multiple value types. At most, one type of value can
// be present, or none at all. Setting a value will clear the other types of
// values so that only one of the 'Is' methods will return true. It is
// possible to clear all values, so that this property is empty.
type ActivityStreamsImagePropertyIterator struct {
	activitystreamsImageMember   vocab.ActivityStreamsImage
	activitystreamsLinkMember    vocab.ActivityStreamsLink
	activitystreamsMentionMember vocab.ActivityStreamsMention
	unknown                      interface{}
	iri                          *url.URL
	alias                        string
	myIdx                        int
	parent                       vocab.ActivityStreamsImageProperty
}

// NewActivityStreamsImagePropertyIterator creates a new ActivityStreamsImage
// property.
func NewActivityStreamsImagePropertyIterator() *ActivityStreamsImagePropertyIterator {
	return &ActivityStreamsImagePropertyIterator{alias: ""}
}

// deserializeActivityStreamsImagePropertyIterator creates an iterator from an
// element that has been unmarshalled from a text or binary format.
func deserializeActivityStreamsImagePropertyIterator(i interface{}, aliasMap map[string]string) (*ActivityStreamsImagePropertyIterator, error) {
	alias := ""
	if a, ok := aliasMap["https://www.w3.org/TR/activitystreams-vocabulary"]; ok {
		alias = a
	}
	if s, ok := i.(string); ok {
		u, err := url.Parse(s)
		// If error exists, don't error out -- skip this and treat as unknown string ([]byte) at worst
		// Also, if no scheme exists, don't treat it as a URL -- net/url is greedy
		if err == nil && len(u.Scheme) > 0 {
			this := &ActivityStreamsImagePropertyIterator{
				alias: alias,
				iri:   u,
			}
			return this, nil
		}
	}
	if m, ok := i.(map[string]interface{}); ok {
		if v, err := mgr.DeserializeImageActivityStreams()(m, aliasMap); err == nil {
			this := &ActivityStreamsImagePropertyIterator{
				activitystreamsImageMember: v,
				alias:                      alias,
			}
			return this, nil
		} else if v, err := mgr.DeserializeLinkActivityStreams()(m, aliasMap); err == nil {
			this := &ActivityStreamsImagePropertyIterator{
				activitystreamsLinkMember: v,
				alias:                     alias,
			}
			return this, nil
		} else if v, err := mgr.DeserializeMentionActivityStreams()(m, aliasMap); err == nil {
			this := &ActivityStreamsImagePropertyIterator{
				activitystreamsMentionMember: v,
				alias:                        alias,
			}
			return this, nil
		}
	}
	this := &ActivityStreamsImagePropertyIterator{
		alias:   alias,
		unknown: i,
	}
	return this, nil
}

// GetActivityStreamsImage returns the value of this property. When
// IsActivityStreamsImage returns false, GetActivityStreamsImage will return
// an arbitrary value.
func (this ActivityStreamsImagePropertyIterator) GetActivityStreamsImage() vocab.ActivityStreamsImage {
	return this.activitystreamsImageMember
}

// GetActivityStreamsLink returns the value of this property. When
// IsActivityStreamsLink returns false, GetActivityStreamsLink will return an
// arbitrary value.
func (this ActivityStreamsImagePropertyIterator) GetActivityStreamsLink() vocab.ActivityStreamsLink {
	return this.activitystreamsLinkMember
}

// GetActivityStreamsMention returns the value of this property. When
// IsActivityStreamsMention returns false, GetActivityStreamsMention will
// return an arbitrary value.
func (this ActivityStreamsImagePropertyIterator) GetActivityStreamsMention() vocab.ActivityStreamsMention {
	return this.activitystreamsMentionMember
}

// GetIRI returns the IRI of this property. When IsIRI returns false, GetIRI will
// return an arbitrary value.
func (this ActivityStreamsImagePropertyIterator) GetIRI() *url.URL {
	return this.iri
}

// GetType returns the value in this property as a Type. Returns nil if the value
// is not an ActivityStreams type, such as an IRI or another value.
func (this ActivityStreamsImagePropertyIterator) GetType() vocab.Type {
	if this.IsActivityStreamsImage() {
		return this.GetActivityStreamsImage()
	}
	if this.IsActivityStreamsLink() {
		return this.GetActivityStreamsLink()
	}
	if this.IsActivityStreamsMention() {
		return this.GetActivityStreamsMention()
	}

	return nil
}

// HasAny returns true if any of the different values is set.
func (this ActivityStreamsImagePropertyIterator) HasAny() bool {
	return this.IsActivityStreamsImage() ||
		this.IsActivityStreamsLink() ||
		this.IsActivityStreamsMention() ||
		this.iri != nil
}

// IsActivityStreamsImage returns true if this property has a type of "Image".
// When true, use the GetActivityStreamsImage and SetActivityStreamsImage
// methods to access and set this property.
func (this ActivityStreamsImagePropertyIterator) IsActivityStreamsImage() bool {
	return this.activitystreamsImageMember != nil
}

// IsActivityStreamsLink returns true if this property has a type of "Link". When
// true, use the GetActivityStreamsLink and SetActivityStreamsLink methods to
// access and set this property.
func (this ActivityStreamsImagePropertyIterator) IsActivityStreamsLink() bool {
	return this.activitystreamsLinkMember != nil
}

// IsActivityStreamsMention returns true if this property has a type of "Mention".
// When true, use the GetActivityStreamsMention and SetActivityStreamsMention
// methods to access and set this property.
func (this ActivityStreamsImagePropertyIterator) IsActivityStreamsMention() bool {
	return this.activitystreamsMentionMember != nil
}

// IsIRI returns true if this property is an IRI. When true, use GetIRI and SetIRI
// to access and set this property
func (this ActivityStreamsImagePropertyIterator) IsIRI() bool {
	return this.iri != nil
}

// JSONLDContext returns the JSONLD URIs required in the context string for this
// property and the specific values that are set. The value in the map is the
// alias used to import the property's value or values.
func (this ActivityStreamsImagePropertyIterator) JSONLDContext() map[string]string {
	m := map[string]string{"https://www.w3.org/TR/activitystreams-vocabulary": this.alias}
	var child map[string]string
	if this.IsActivityStreamsImage() {
		child = this.GetActivityStreamsImage().JSONLDContext()
	} else if this.IsActivityStreamsLink() {
		child = this.GetActivityStreamsLink().JSONLDContext()
	} else if this.IsActivityStreamsMention() {
		child = this.GetActivityStreamsMention().JSONLDContext()
	}
	/*
	   Since the literal maps in this function are determined at
	   code-generation time, this loop should not overwrite an existing key with a
	   new value.
	*/
	for k, v := range child {
		m[k] = v
	}
	return m
}

// KindIndex computes an arbitrary value for indexing this kind of value. This is
// a leaky API detail only for folks looking to replace the go-fed
// implementation. Applications should not use this method.
func (this ActivityStreamsImagePropertyIterator) KindIndex() int {
	if this.IsActivityStreamsImage() {
		return 0
	}
	if this.IsActivityStreamsLink() {
		return 1
	}
	if this.IsActivityStreamsMention() {
		return 2
	}
	if this.IsIRI() {
		return -2
	}
	return -1
}

// LessThan compares two instances of this property with an arbitrary but stable
// comparison. Applications should not use this because it is only meant to
// help alternative implementations to go-fed to be able to normalize
// nonfunctional properties.
func (this ActivityStreamsImagePropertyIterator) LessThan(o vocab.ActivityStreamsImagePropertyIterator) bool {
	idx1 := this.KindIndex()
	idx2 := o.KindIndex()
	if idx1 < idx2 {
		return true
	} else if idx1 > idx2 {
		return false
	} else if this.IsActivityStreamsImage() {
		return this.GetActivityStreamsImage().LessThan(o.GetActivityStreamsImage())
	} else if this.IsActivityStreamsLink() {
		return this.GetActivityStreamsLink().LessThan(o.GetActivityStreamsLink())
	} else if this.IsActivityStreamsMention() {
		return this.GetActivityStreamsMention().LessThan(o.GetActivityStreamsMention())
	} else if this.IsIRI() {
		return this.iri.String() < o.GetIRI().String()
	}
	return false
}

// Name returns the name of this property: "ActivityStreamsImage".
func (this ActivityStreamsImagePropertyIterator) Name() string {
	return "ActivityStreamsImage"
}

// Next returns the next iterator, or nil if there is no next iterator.
func (this ActivityStreamsImagePropertyIterator) Next() vocab.ActivityStreamsImagePropertyIterator {
	if this.myIdx+1 >= this.parent.Len() {
		return nil
	} else {
		return this.parent.At(this.myIdx + 1)
	}
}

// Prev returns the previous iterator, or nil if there is no previous iterator.
func (this ActivityStreamsImagePropertyIterator) Prev() vocab.ActivityStreamsImagePropertyIterator {
	if this.myIdx-1 < 0 {
		return nil
	} else {
		return this.parent.At(this.myIdx - 1)
	}
}

// SetActivityStreamsImage sets the value of this property. Calling
// IsActivityStreamsImage afterwards returns true.
func (this *ActivityStreamsImagePropertyIterator) SetActivityStreamsImage(v vocab.ActivityStreamsImage) {
	this.clear()
	this.activitystreamsImageMember = v
}

// SetActivityStreamsLink sets the value of this property. Calling
// IsActivityStreamsLink afterwards returns true.
func (this *ActivityStreamsImagePropertyIterator) SetActivityStreamsLink(v vocab.ActivityStreamsLink) {
	this.clear()
	this.activitystreamsLinkMember = v
}

// SetActivityStreamsMention sets the value of this property. Calling
// IsActivityStreamsMention afterwards returns true.
func (this *ActivityStreamsImagePropertyIterator) SetActivityStreamsMention(v vocab.ActivityStreamsMention) {
	this.clear()
	this.activitystreamsMentionMember = v
}

// SetIRI sets the value of this property. Calling IsIRI afterwards returns true.
func (this *ActivityStreamsImagePropertyIterator) SetIRI(v *url.URL) {
	this.clear()
	this.iri = v
}

// SetType attempts to set the property for the arbitrary type. Returns an error
// if it is not a valid type to set on this property.
func (this *ActivityStreamsImagePropertyIterator) SetType(t vocab.Type) error {
	if v, ok := t.(vocab.ActivityStreamsImage); ok {
		this.SetActivityStreamsImage(v)
		return nil
	}
	if v, ok := t.(vocab.ActivityStreamsLink); ok {
		this.SetActivityStreamsLink(v)
		return nil
	}
	if v, ok := t.(vocab.ActivityStreamsMention); ok {
		this.SetActivityStreamsMention(v)
		return nil
	}

	return fmt.Errorf("illegal type to set on ActivityStreamsImage property: %T", t)
}

// clear ensures no value of this property is set. Calling HasAny or any of the
// 'Is' methods afterwards will return false.
func (this *ActivityStreamsImagePropertyIterator) clear() {
	this.activitystreamsImageMember = nil
	this.activitystreamsLinkMember = nil
	this.activitystreamsMentionMember = nil
	this.unknown = nil
	this.iri = nil
}

// serialize converts this into an interface representation suitable for
// marshalling into a text or binary format. Applications should not need this
// function as most typical use cases serialize types instead of individual
// properties. It is exposed for alternatives to go-fed implementations to use.
func (this ActivityStreamsImagePropertyIterator) serialize() (interface{}, error) {
	if this.IsActivityStreamsImage() {
		return this.GetActivityStreamsImage().Serialize()
	} else if this.IsActivityStreamsLink() {
		return this.GetActivityStreamsLink().Serialize()
	} else if this.IsActivityStreamsMention() {
		return this.GetActivityStreamsMention().Serialize()
	} else if this.IsIRI() {
		return this.iri.String(), nil
	}
	return this.unknown, nil
}

// ActivityStreamsImageProperty is the non-functional property "image". It is
// permitted to have one or more values, and of different value types.
type ActivityStreamsImageProperty struct {
	properties []*ActivityStreamsImagePropertyIterator
	alias      string
}

// DeserializeImageProperty creates a "image" property from an interface
// representation that has been unmarshalled from a text or binary format.
func DeserializeImageProperty(m map[string]interface{}, aliasMap map[string]string) (vocab.ActivityStreamsImageProperty, error) {
	alias := ""
	if a, ok := aliasMap["https://www.w3.org/TR/activitystreams-vocabulary"]; ok {
		alias = a
	}
	propName := "image"
	if len(alias) > 0 {
		propName = fmt.Sprintf("%s:%s", alias, "image")
	}
	i, ok := m[propName]

	if ok {
		this := &ActivityStreamsImageProperty{
			alias:      alias,
			properties: []*ActivityStreamsImagePropertyIterator{},
		}
		if list, ok := i.([]interface{}); ok {
			for _, iterator := range list {
				if p, err := deserializeActivityStreamsImagePropertyIterator(iterator, aliasMap); err != nil {
					return this, err
				} else if p != nil {
					this.properties = append(this.properties, p)
				}
			}
		} else {
			if p, err := deserializeActivityStreamsImagePropertyIterator(i, aliasMap); err != nil {
				return this, err
			} else if p != nil {
				this.properties = append(this.properties, p)
			}
		}
		// Set up the properties for iteration.
		for idx, ele := range this.properties {
			ele.parent = this
			ele.myIdx = idx
		}
		return this, nil
	}
	return nil, nil
}

// NewActivityStreamsImageProperty creates a new image property.
func NewActivityStreamsImageProperty() *ActivityStreamsImageProperty {
	return &ActivityStreamsImageProperty{alias: ""}
}

// AppendActivityStreamsImage appends a Image value to the back of a list of the
// property "image". Invalidates iterators that are traversing using Prev.
func (this *ActivityStreamsImageProperty) AppendActivityStreamsImage(v vocab.ActivityStreamsImage) {
	this.properties = append(this.properties, &ActivityStreamsImagePropertyIterator{
		activitystreamsImageMember: v,
		alias:                      this.alias,
		myIdx:                      this.Len(),
		parent:                     this,
	})
}

// AppendActivityStreamsLink appends a Link value to the back of a list of the
// property "image". Invalidates iterators that are traversing using Prev.
func (this *ActivityStreamsImageProperty) AppendActivityStreamsLink(v vocab.ActivityStreamsLink) {
	this.properties = append(this.properties, &ActivityStreamsImagePropertyIterator{
		activitystreamsLinkMember: v,
		alias:                     this.alias,
		myIdx:                     this.Len(),
		parent:                    this,
	})
}

// AppendActivityStreamsMention appends a Mention value to the back of a list of
// the property "image". Invalidates iterators that are traversing using Prev.
func (this *ActivityStreamsImageProperty) AppendActivityStreamsMention(v vocab.ActivityStreamsMention) {
	this.properties = append(this.properties, &ActivityStreamsImagePropertyIterator{
		activitystreamsMentionMember: v,
		alias:                        this.alias,
		myIdx:                        this.Len(),
		parent:                       this,
	})
}

// AppendIRI appends an IRI value to the back of a list of the property "image"
func (this *ActivityStreamsImageProperty) AppendIRI(v *url.URL) {
	this.properties = append(this.properties, &ActivityStreamsImagePropertyIterator{
		alias:  this.alias,
		iri:    v,
		myIdx:  this.Len(),
		parent: this,
	})
}

// PrependType prepends an arbitrary type value to the front of a list of the
// property "image". Invalidates iterators that are traversing using Prev.
// Returns an error if the type is not a valid one to set for this property.
func (this *ActivityStreamsImageProperty) AppendType(t vocab.Type) error {
	n := &ActivityStreamsImagePropertyIterator{
		alias:  this.alias,
		myIdx:  this.Len(),
		parent: this,
	}
	if err := n.SetType(t); err != nil {
		return err
	}
	this.properties = append(this.properties, n)
	return nil
}

// At returns the property value for the specified index. Panics if the index is
// out of bounds.
func (this ActivityStreamsImageProperty) At(index int) vocab.ActivityStreamsImagePropertyIterator {
	return this.properties[index]
}

// Begin returns the first iterator, or nil if empty. Can be used with the
// iterator's Next method and this property's End method to iterate from front
// to back through all values.
func (this ActivityStreamsImageProperty) Begin() vocab.ActivityStreamsImagePropertyIterator {
	if this.Empty() {
		return nil
	} else {
		return this.properties[0]
	}
}

// Empty returns returns true if there are no elements.
func (this ActivityStreamsImageProperty) Empty() bool {
	return this.Len() == 0
}

// End returns beyond-the-last iterator, which is nil. Can be used with the
// iterator's Next method and this property's Begin method to iterate from
// front to back through all values.
func (this ActivityStreamsImageProperty) End() vocab.ActivityStreamsImagePropertyIterator {
	return nil
}

// InsertActivityStreamsImage inserts a Image value at the specified index for a
// property "image". Existing elements at that index and higher are shifted
// back once. Invalidates all iterators.
func (this *ActivityStreamsImageProperty) InsertActivityStreamsImage(idx int, v vocab.ActivityStreamsImage) {
	this.properties = append(this.properties, nil)
	copy(this.properties[idx+1:], this.properties[idx:])
	this.properties[idx] = &ActivityStreamsImagePropertyIterator{
		activitystreamsImageMember: v,
		alias:                      this.alias,
		myIdx:                      idx,
		parent:                     this,
	}
	for i := idx; i < this.Len(); i++ {
		(this.properties)[i].myIdx = i
	}
}

// InsertActivityStreamsLink inserts a Link value at the specified index for a
// property "image". Existing elements at that index and higher are shifted
// back once. Invalidates all iterators.
func (this *ActivityStreamsImageProperty) InsertActivityStreamsLink(idx int, v vocab.ActivityStreamsLink) {
	this.properties = append(this.properties, nil)
	copy(this.properties[idx+1:], this.properties[idx:])
	this.properties[idx] = &ActivityStreamsImagePropertyIterator{
		activitystreamsLinkMember: v,
		alias:                     this.alias,
		myIdx:                     idx,
		parent:                    this,
	}
	for i := idx; i < this.Len(); i++ {
		(this.properties)[i].myIdx = i
	}
}

// InsertActivityStreamsMention inserts a Mention value at the specified index for
// a property "image". Existing elements at that index and higher are shifted
// back once. Invalidates all iterators.
func (this *ActivityStreamsImageProperty) InsertActivityStreamsMention(idx int, v vocab.ActivityStreamsMention) {
	this.properties = append(this.properties, nil)
	copy(this.properties[idx+1:], this.properties[idx:])
	this.properties[idx] = &ActivityStreamsImagePropertyIterator{
		activitystreamsMentionMember: v,
		alias:                        this.alias,
		myIdx:                        idx,
		parent:                       this,
	}
	for i := idx; i < this.Len(); i++ {
		(this.properties)[i].myIdx = i
	}
}

// Insert inserts an IRI value at the specified index for a property "image".
// Existing elements at that index and higher are shifted back once.
// Invalidates all iterators.
func (this *ActivityStreamsImageProperty) InsertIRI(idx int, v *url.URL) {
	this.properties = append(this.properties, nil)
	copy(this.properties[idx+1:], this.properties[idx:])
	this.properties[idx] = &ActivityStreamsImagePropertyIterator{
		alias:  this.alias,
		iri:    v,
		myIdx:  idx,
		parent: this,
	}
	for i := idx; i < this.Len(); i++ {
		(this.properties)[i].myIdx = i
	}
}

// PrependType prepends an arbitrary type value to the front of a list of the
// property "image". Invalidates all iterators. Returns an error if the type
// is not a valid one to set for this property.
func (this *ActivityStreamsImageProperty) InsertType(idx int, t vocab.Type) error {
	n := &ActivityStreamsImagePropertyIterator{
		alias:  this.alias,
		myIdx:  idx,
		parent: this,
	}
	if err := n.SetType(t); err != nil {
		return err
	}
	this.properties = append(this.properties, nil)
	copy(this.properties[idx+1:], this.properties[idx:])
	this.properties[idx] = n
	for i := idx; i < this.Len(); i++ {
		(this.properties)[i].myIdx = i
	}
	return nil
}

// JSONLDContext returns the JSONLD URIs required in the context string for this
// property and the specific values that are set. The value in the map is the
// alias used to import the property's value or values.
func (this ActivityStreamsImageProperty) JSONLDContext() map[string]string {
	m := map[string]string{"https://www.w3.org/TR/activitystreams-vocabulary": this.alias}
	for _, elem := range this.properties {
		child := elem.JSONLDContext()
		/*
		   Since the literal maps in this function are determined at
		   code-generation time, this loop should not overwrite an existing key with a
		   new value.
		*/
		for k, v := range child {
			m[k] = v
		}
	}
	return m
}

// KindIndex computes an arbitrary value for indexing this kind of value. This is
// a leaky API method specifically needed only for alternate implementations
// for go-fed. Applications should not use this method. Panics if the index is
// out of bounds.
func (this ActivityStreamsImageProperty) KindIndex(idx int) int {
	return this.properties[idx].KindIndex()
}

// Len returns the number of values that exist for the "image" property.
func (this ActivityStreamsImageProperty) Len() (length int) {
	return len(this.properties)
}

// Less computes whether another property is less than this one. Mixing types
// results in a consistent but arbitrary ordering
func (this ActivityStreamsImageProperty) Less(i, j int) bool {
	idx1 := this.KindIndex(i)
	idx2 := this.KindIndex(j)
	if idx1 < idx2 {
		return true
	} else if idx1 == idx2 {
		if idx1 == 0 {
			lhs := this.properties[i].GetActivityStreamsImage()
			rhs := this.properties[j].GetActivityStreamsImage()
			return lhs.LessThan(rhs)
		} else if idx1 == 1 {
			lhs := this.properties[i].GetActivityStreamsLink()
			rhs := this.properties[j].GetActivityStreamsLink()
			return lhs.LessThan(rhs)
		} else if idx1 == 2 {
			lhs := this.properties[i].GetActivityStreamsMention()
			rhs := this.properties[j].GetActivityStreamsMention()
			return lhs.LessThan(rhs)
		} else if idx1 == -2 {
			lhs := this.properties[i].GetIRI()
			rhs := this.properties[j].GetIRI()
			return lhs.String() < rhs.String()
		}
	}
	return false
}

// LessThan compares two instances of this property with an arbitrary but stable
// comparison. Applications should not use this because it is only meant to
// help alternative implementations to go-fed to be able to normalize
// nonfunctional properties.
func (this ActivityStreamsImageProperty) LessThan(o vocab.ActivityStreamsImageProperty) bool {
	l1 := this.Len()
	l2 := o.Len()
	l := l1
	if l2 < l1 {
		l = l2
	}
	for i := 0; i < l; i++ {
		if this.properties[i].LessThan(o.At(i)) {
			return true
		} else if o.At(i).LessThan(this.properties[i]) {
			return false
		}
	}
	return l1 < l2
}

// Name returns the name of this property: "image".
func (this ActivityStreamsImageProperty) Name() string {
	return "image"
}

// PrependActivityStreamsImage prepends a Image value to the front of a list of
// the property "image". Invalidates all iterators.
func (this *ActivityStreamsImageProperty) PrependActivityStreamsImage(v vocab.ActivityStreamsImage) {
	this.properties = append([]*ActivityStreamsImagePropertyIterator{{
		activitystreamsImageMember: v,
		alias:                      this.alias,
		myIdx:                      0,
		parent:                     this,
	}}, this.properties...)
	for i := 1; i < this.Len(); i++ {
		(this.properties)[i].myIdx = i
	}
}

// PrependActivityStreamsLink prepends a Link value to the front of a list of the
// property "image". Invalidates all iterators.
func (this *ActivityStreamsImageProperty) PrependActivityStreamsLink(v vocab.ActivityStreamsLink) {
	this.properties = append([]*ActivityStreamsImagePropertyIterator{{
		activitystreamsLinkMember: v,
		alias:                     this.alias,
		myIdx:                     0,
		parent:                    this,
	}}, this.properties...)
	for i := 1; i < this.Len(); i++ {
		(this.properties)[i].myIdx = i
	}
}

// PrependActivityStreamsMention prepends a Mention value to the front of a list
// of the property "image". Invalidates all iterators.
func (this *ActivityStreamsImageProperty) PrependActivityStreamsMention(v vocab.ActivityStreamsMention) {
	this.properties = append([]*ActivityStreamsImagePropertyIterator{{
		activitystreamsMentionMember: v,
		alias:                        this.alias,
		myIdx:                        0,
		parent:                       this,
	}}, this.properties...)
	for i := 1; i < this.Len(); i++ {
		(this.properties)[i].myIdx = i
	}
}

// PrependIRI prepends an IRI value to the front of a list of the property "image".
func (this *ActivityStreamsImageProperty) PrependIRI(v *url.URL) {
	this.properties = append([]*ActivityStreamsImagePropertyIterator{{
		alias:  this.alias,
		iri:    v,
		myIdx:  0,
		parent: this,
	}}, this.properties...)
	for i := 1; i < this.Len(); i++ {
		(this.properties)[i].myIdx = i
	}
}

// PrependType prepends an arbitrary type value to the front of a list of the
// property "image". Invalidates all iterators. Returns an error if the type
// is not a valid one to set for this property.
func (this *ActivityStreamsImageProperty) PrependType(t vocab.Type) error {
	n := &ActivityStreamsImagePropertyIterator{
		alias:  this.alias,
		myIdx:  0,
		parent: this,
	}
	if err := n.SetType(t); err != nil {
		return err
	}
	this.properties = append([]*ActivityStreamsImagePropertyIterator{n}, this.properties...)
	for i := 1; i < this.Len(); i++ {
		(this.properties)[i].myIdx = i
	}
	return nil
}

// Remove deletes an element at the specified index from a list of the property
// "image", regardless of its type. Panics if the index is out of bounds.
// Invalidates all iterators.
func (this *ActivityStreamsImageProperty) Remove(idx int) {
	(this.properties)[idx].parent = nil
	copy((this.properties)[idx:], (this.properties)[idx+1:])
	(this.properties)[len(this.properties)-1] = &ActivityStreamsImagePropertyIterator{}
	this.properties = (this.properties)[:len(this.properties)-1]
	for i := idx; i < this.Len(); i++ {
		(this.properties)[i].myIdx = i
	}
}

// Serialize converts this into an interface representation suitable for
// marshalling into a text or binary format. Applications should not need this
// function as most typical use cases serialize types instead of individual
// properties. It is exposed for alternatives to go-fed implementations to use.
func (this ActivityStreamsImageProperty) Serialize() (interface{}, error) {
	s := make([]interface{}, 0, len(this.properties))
	for _, iterator := range this.properties {
		if b, err := iterator.serialize(); err != nil {
			return s, err
		} else {
			s = append(s, b)
		}
	}
	// Shortcut: if serializing one value, don't return an array -- pretty sure other Fediverse software would choke on a "type" value with array, for example.
	if len(s) == 1 {
		return s[0], nil
	}
	return s, nil
}

// SetActivityStreamsImage sets a Image value to be at the specified index for the
// property "image". Panics if the index is out of bounds. Invalidates all
// iterators.
func (this *ActivityStreamsImageProperty) SetActivityStreamsImage(idx int, v vocab.ActivityStreamsImage) {
	(this.properties)[idx].parent = nil
	(this.properties)[idx] = &ActivityStreamsImagePropertyIterator{
		activitystreamsImageMember: v,
		alias:                      this.alias,
		myIdx:                      idx,
		parent:                     this,
	}
}

// SetActivityStreamsLink sets a Link value to be at the specified index for the
// property "image". Panics if the index is out of bounds. Invalidates all
// iterators.
func (this *ActivityStreamsImageProperty) SetActivityStreamsLink(idx int, v vocab.ActivityStreamsLink) {
	(this.properties)[idx].parent = nil
	(this.properties)[idx] = &ActivityStreamsImagePropertyIterator{
		activitystreamsLinkMember: v,
		alias:                     this.alias,
		myIdx:                     idx,
		parent:                    this,
	}
}

// SetActivityStreamsMention sets a Mention value to be at the specified index for
// the property "image". Panics if the index is out of bounds. Invalidates all
// iterators.
func (this *ActivityStreamsImageProperty) SetActivityStreamsMention(idx int, v vocab.ActivityStreamsMention) {
	(this.properties)[idx].parent = nil
	(this.properties)[idx] = &ActivityStreamsImagePropertyIterator{
		activitystreamsMentionMember: v,
		alias:                        this.alias,
		myIdx:                        idx,
		parent:                       this,
	}
}

// SetIRI sets an IRI value to be at the specified index for the property "image".
// Panics if the index is out of bounds.
func (this *ActivityStreamsImageProperty) SetIRI(idx int, v *url.URL) {
	(this.properties)[idx].parent = nil
	(this.properties)[idx] = &ActivityStreamsImagePropertyIterator{
		alias:  this.alias,
		iri:    v,
		myIdx:  idx,
		parent: this,
	}
}

// SetType sets an arbitrary type value to the specified index of the property
// "image". Invalidates all iterators. Returns an error if the type is not a
// valid one to set for this property. Panics if the index is out of bounds.
func (this *ActivityStreamsImageProperty) SetType(idx int, t vocab.Type) error {
	n := &ActivityStreamsImagePropertyIterator{
		alias:  this.alias,
		myIdx:  idx,
		parent: this,
	}
	if err := n.SetType(t); err != nil {
		return err
	}
	(this.properties)[idx] = n
	return nil
}

// Swap swaps the location of values at two indices for the "image" property.
func (this ActivityStreamsImageProperty) Swap(i, j int) {
	this.properties[i], this.properties[j] = this.properties[j], this.properties[i]
}
