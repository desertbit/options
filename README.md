Suggests a design guideline for Config/Option structs.

```go
type Options struct {
    NumWorkers int `yaml:"num-workers,omitempty"`
    Name string `yaml:"name,omitempty"`
    ConfThreshold int `yaml:"conf-threshold,omitempty"`
}

func DefaultOptions() Options {
    return Options{
        NumWorkers: 8,
        Name: "testName",
        ConfThreshold: 0.85,
    }
}

func (o Options) Save(path string) error {
    // Do not serialize any values that match the default values.
    // This way, only changes are reflected in the options file.
    // The Options struct itself is not altered, since the receiver
    // is not a pointer.
    err = options.StripDefaults(&o, DefaultOptions())
    if err != nil {
        return err
    }

    // ... yaml marshal and write to file ...
}
```

The main concept is **simplicity**. One simple Options struct combined with one constructor that returns the default settings.
By not serializing any default values (`StripDefaults`), we keep the written file tidy and only list **changes to the defaults**.

In case a **slice** or **map** is used in the config, whose elements are structs again, it is not possible to state real default values.
Instead, one should iterate over the options after parsing and set default values on each of the elements.

```go
type Options struct {
    NumWorkers int `yaml:"num-workers,omitempty"`
    Name string `yaml:"name,omitempty"`
    ConfThreshold int `yaml:"conf-threshold,omitempty"`
    Sl []Nested `yaml:"sl,omitempty"`
}

type Nested struct {
    Test int `yaml:"test,omitempty"`
}

func DefaultNestedOptions() Nested {
    return Nested{Test: 5}
}

func ParseOptionsFromFile(path string) (o Options, err error) {
	o = DefaultOptions(path)

	// Read data from file.
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return
	}

	// Unmarshal the data.
	err = yaml.UnmarshalStrict(data, &o)
	if err != nil {
		return
	}

	// Iterate over the nested options and set default values.
	dn := DefaultNestedOptions()
	for i := range o.Sl {
	    err = options.SetDefaults(o.Sl[i], dn)
	    if err != nil {
	        return
	    }
	}

	return
}
```
