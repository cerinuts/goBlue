/*
Copyright (c) 2018 ceriath
This Package is part of the "goBlue"-Library
It is licensed under the MIT License
*/

//Package network offers various network tools
package network

//AppName is the name of the application
const AppName string = "goBlue/net"

//VersionMajor 0 means in development, >1 ensures compatibility with each minor version, but breakes with new major version
const VersionMajor string = "0"

//VersionMinor introduces changes that require a new version number. If the major version is 0, they are likely to break compatibility
const VersionMinor string = "3"

//VersionBuild is the type of this release. s(table), b(eta), d(evelopment), n(ightly)
const VersionBuild string = "s"

//FullVersion contains the full name and version of this package in a printable string
const FullVersion string = AppName + VersionMajor + "." + VersionMinor + VersionBuild
