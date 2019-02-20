===============================================================================
goluamapper: maps an Azure/golua table to a Go struct
===============================================================================

.. image:: https://godoc.org/github.com/jdolitsky/gluamapper?status.svg
    :target: http://godoc.org/github.com/jdolitsky/gluamapper

|

goluamapper provides an easy way to map Azure/golua tables to Go structs.

goluamapper converts an Azure/golua table to ``map[string]interface{}`` , and then converts it to a Go struct using `mapstructure <https://github.com/mitchellh/mapstructure/>`_ .

----------------------------------------------------------------
API
----------------------------------------------------------------
See `Go doc <http://godoc.org/github.com/jdolitsky/gluamapper>`_ .

----------------------------------------------------------------
Usage
----------------------------------------------------------------

.. code-block:: go

    type Role struct {
        Name string
    }

    type Person struct {
        Name      string
        Age       int
        WorkPlace string
        Role      []*Role
    }

    L := lua.NewState()
    if err := L.DoString(`
    person = {
      name = "Michel",
      age  = "31", -- weakly input
      work_place = "San Jose",
      role = {
        {
          name = "Administrator"
        },
        {
          name = "Operator"
        }
      }
    }
    `); err != nil {
        panic(err)
    }
    var person Person
    if err := goluamapper.Map(L.GetGlobal("person").(*lua.LTable), &person); err != nil {
        panic(err)
    }
    fmt.Printf("%s %d", person.Name, person.Age)

----------------------------------------------------------------
License
----------------------------------------------------------------
MIT

----------------------------------------------------------------
Original Author
----------------------------------------------------------------
Yusuke Inuzuka

Source: https://github.com/yuin/gluamapper
