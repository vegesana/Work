package rest

const d10VDoc = `
<!DOCTYPE html>
<html>
    <head><title>{{.Title}}</title></head>
    <body>
        <h1>{{.Title}}</h1>
		<h5>LSB/RightMostBit : Lif-0</h5>
		<table style="width:26%" border="3">
        <tr>
            	<th> Server </th>
            	<th> FLOOD BITMAP</th>
        </tr>
			{{template "List" .D10VlanSlice}}
        </table>
    </body>
</html>
`

const d10VDocList = `
<ul >
    {{range .}}
    <tr>
        <td>{{.Server}}</td>
	        {{template "List1" .D10Vlan}}
    </tr>
    {{end}}
</ul>
`
const d10VDocList1 = `
<ul >
	<td>
		<table style="width:25%" border="1">
		<tr>
			<th> VLANID </th>
			<th> LIF BITMAP</th>
		</tr>
    {{range .}}
		<tr>
        	<td>{{.Vlan}}</td>
        	<td>{{.LifBitmap}}</td>
		</tr>
    {{end}}
	</td>
	</table>
</ul>
`

const lifDoc = `
<!DOCTYPE html>
<html>
    <head><title>{{.Title}}</title></head>
    <body>
        <h1>{{.Title}}</h1>
		<table style="width:23%" border="3">
        <tr>
            	<th> Server </th>
            	<th> LIFID,VLANID,PHYIF </th>
        </tr>
	        {{template "List" .D10lifSlice}}
        </table>
    </body>
</html>
`

const lifDocList = `
<ul >
    {{range .}}
    <tr>
        <td>{{.Server}}</td>
	        {{template "List1" .LifSlice}}
    </tr>
    {{end}}
</ul>
`
const lifDocList1 = `
<ul >
	<td>
		<table style="width:24%" border="1">
		<tr>
			<th> LIFID </th>
			<th> VLANID</th>
			<th> PHYIF INFO </th>
		</tr>
    {{range .}}
		<tr>
        	<td>{{.LifId}}</td>
        	<td>{{.Vlan}}</td>
        	<td>{{.Pif}}</td>
		</tr>
    {{end}}
	</td>
	</table>
</ul>
`

const ctrlDocList = `
<ul >
    {{range .}}
    <tr>
        <td>{{.Servername}}</td>
        <td>{{.Id}}</td>
        <td>{{.Type}}</td>
        <td>{{.LocalMac}}</td>
        <td>{{.RemoteMac}}</td>
        <td>{{.PairedCtrlId}}</td>
        <td>{{.Cookie}}</td>
    </tr>
    {{end}}
</ul>
`
const ctrlDoc = `
<!DOCTYPE html>
<html>
    <head><title>{{.Title}}</title></head>
    <body>
        <h1>{{.Title}}</h1>
		<table style="width:50%" border="3">
        <tr>
            	<th> Server </th>
            	<th> Ctrl ID </th>
            	<th> Type </th>
            	<th> Local Mac </th>
            	<th> Remote Mac </th>
            	<th> Remote Ctrl Id </th>
            	<th> Cookie </th>
        </tr>
	        {{template "List" .CtrlInfor}}
        </table>
    </body>
</html>
`
const nodeDocList = `
<ul >
    {{range .}}
    <tr>
        <td>{{.HostName}}</td>
        <td>{{.NodeHealth}}</td>
        <td>{{.K8sHealth}}</td>
        <td>{{.Zone}}</td>
        <td>{{.SIP1}}</td>
        <td>{{.SIP3}}</td>
        <td>{{.Svlan}}</td>
        <td>{{.GMAC}}</td>
        <td>{{.Mode}}</td>
    </tr>
    {{end}}
</ul>
`
const nodeDoc = `
<!DOCTYPE html>
<html>
    <head><title>{{.Title}}</title></head>
    <body>
        <h1>{{.Title}}</h1>
		<table style="width:50%" border="3">
        <tr>
            	<th> HostName </th>
            	<th> Node Health </th>
            	<th> K8s Health </th>
            	<th> Zone </th>
            	<th> Storage IP 1 </th>
            	<th> Storage IP 3 </th>
            	<th> Storage VLAN </th>
            	<th> GWMAC </th>
            	<th> Mode </th>
        </tr>
	        {{template "List" .NodeInfor}}
        </table>
    </body>
</html>
`

const networkDocList = ` <ul >
    {{range .}}
    <tr>
        <td>{{.Name}}</td>
        <td>{{.Type}}</td>
        <td>{{.SSubnet}}</td>
        <td>{{.SVlan}}</td>
        <td>{{.SNumAddr}}</td>
        <td>{{.SUsedAddr}}</td>
        <td>{{.GatewayMac}}</td>
        <td>{{.Zone}}</td>
    </tr>
    {{end}}
</ul>
`
const networkDoc = `
<!DOCTYPE html>
<html>
    <head><title>{{.Title}}</title></head>
    <body>
        <h1>{{.Title}}</h1>
		<table style="width:50%" border="3">
        <tr>
            	<th> Network Name </th>
            	<th> Type </th>
            	<th> Subnet </th>
            	<th> SVlan </th>
            	<th> TotalAddrCnt </th>
            	<th> UsedAddrCnt </th>
            	<th> GWMac </th>
            	<th> Zone </th>
        </tr>
	        {{template "List" .NetworkInfor}}
        </table>
    </body>
</html>
`

const intfHelpDoc = `
<!DOCTYPE html>
<html>
    <head><title>{{.Title}}</title></head>
    <body>
        <h1>{{.Title}}</h1>
		<table style="width:15%" border="3">
        <tr>
            	<th> intf_type </th>
            	<th> Description </th>
        </tr>
	        {{template "List" .IntfHelpInfo}}
        </table>
    </body>
</html>
`

const intfHelpDocList = `
<ul >
    {{range .}}
    <tr>
        <td>{{.Help}}</td>
        <td>{{.Desc}}</td>
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
