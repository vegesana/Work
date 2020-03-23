package rest

const pclDocList1 = `
<ul >
	<td>
		<table style="width:140%" border="0">
    {{range .}}
		<tr>
        	<td>{{.Rule}}</td>
		</tr>
    {{end}}
	</td>
	</table>
</ul>
`

const pclDocList = `
<ul >
    {{range .}}
    <tr>
        <td>{{.Server}}</td>
	        {{template "List1" .RuleList}}
    </tr>
    {{end}}
</ul>
`
const pclDoc = `
<!DOCTYPE html>
<html>
    <head><title>{{.Title}}</title></head>
    <body>
        <h1>{{.Title}}</h1>
		<table style="width:80%" border="3">
        <tr>
            	<th> Server </th>
            	<th> Rule List </th>
        </tr>
	        {{template "List" .PclInfo}}
        </table>
    </body>
</html>
`

const pinDocList1 = `
<ul >
	<td>
		<table style="width:5%" border="0">
    {{range .}}
		<tr>
        	<td>{{.SPort}}</td>
        	<td> - </td>
        	<td>{{.DPort}}</td>
		</tr>
    {{end}}
	</td>
	</table>
</ul>
`

const pinDocList = `
<ul >
    {{range .}}
    <tr>
        <td>{{.Server}}</td>
	        {{template "List1" .PortSlice}}
    </tr>
    {{end}}
</ul>
`
const pinDoc = `
<!DOCTYPE html>
<html>
    <head><title>{{.Title}}</title></head>
    <body>
        <h1>{{.Title}}</h1>
		<table style="width:14%" border="3">
        <tr>
            	<th> Server </th>
            	<th> SrcPort - DestPort</th>
        </tr>
	        {{template "List" .PinInfo}}
        </table>
    </body>
</html>
`

const macDocList1 = `
<ul >
	<td>
		<table border="1">
    {{range .}}
		<tr>
        	<td>{{.Server}}</td>
        	<td>{{.Port}}</td>
		</tr>
    {{end}}
	</td>
	</table>
</ul>
`

const macDocList = `
<ul >
    {{range .}}
    <tr>
        <td>{{.Mac}}</td>
        <td>{{.Vlan}}</td>
	        {{template "List1" .ServerPortSlice}}
    </tr>
    {{end}}
</ul>
`

const macDoc = `
<!DOCTYPE html>
<html>
    <head><title>{{.Title}}</title></head>
    <body>
        <h1>{{.Title}}</h1>
		<table style="width:20%" border="3">
        <tr>
            	<th> MAC</th>
            	<th> Vlan</th>
            	<th> Server-Port</th>
        </tr>
	        {{template "List" .MacInfo}}
        </table>
    </body>
</html>
`
const errDocList = `
<ul >
    {{range .}}
    <tr>
        <td>{{.ServerName}}</td>
        <td>{{.MyErr}}</td>
    </tr>
    {{end}}
</ul>
`

const errDoc = `
<!DOCTYPE html>
<html>
    <head><title>{{.Title}}</title></head>
    <body>
        <h1>{{.Title}}</h1>
        <table style="width:40%" border="1">
        <tr>
            <th> ServerName </th>
            <th> Error </th>
        </tr>
        {{template "List" .ErrInfo}}
        </table>
    </body>
</html>
`
const sysDocList = `
<ul >
    {{range .}}
    <tr>
        <td>{{.ServerName}}</td>
        <td>{{.BoardInfo}}</td>
        <td>{{.ProductId}}</td>
    </tr>
    {{end}}
</ul>
`

const sysDoc = `
<!DOCTYPE html>
<html>
    <head><title>{{.Title}}</title></head>
    <body>
        <h1>{{.Title}}</h1>
        <table style="width:20%" border="1">
        <tr>
            <th> ServerName </th>
            <th> BoardInfo</th>
            <th> ProductId</th>
        </tr>
        {{template "List" .SysInfo}}
        </table>
    </body>
</html>
`