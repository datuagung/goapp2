

package main

import (
    "github.com/spf13/viper"
)

func init() {
    viper.SetConfigFile(`config.json`)
    err := viper.ReadInConfig()

    if err != nil {
        panic(err)
    }

}

func main() {

    dbHost := viper.GetString(`database.host`)
    dbPort := viper.GetString(`database.port`)
    dbUser := viper.GetString(`database.user`)
    dbPass := viper.GetString(`database.pass`)
    dbName := viper.GetString(`database.name`)
    dbaddress := viper.GetString(`server.address`)

    a := App{}
    a.Initialize(dbUser, dbPass, dbHost, dbPort, dbName)
    a.Run(dbaddress)
}