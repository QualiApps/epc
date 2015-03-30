# EPC

EPC is a Go client library for accessing the Epam Private Cloud.

## Usage

```go
import "github.com/qualiapps/epc"
```

Create a new EPC client, then use the exposed services to access different parts of the EPC API.

## Examples


To create a new Instance:

```go
func GetClient() *epc.EPC {
    auth := epc.GetAuth(AccessUser, AccessToken)
    return epc.NewEPC(auth, ProjectID, Zone, Image, Shape)
}

//Create key pairs
...

// Run instance
instance, err := GetClient().CreateInstance(KeyName)
if err != nil {
    return err
}
```