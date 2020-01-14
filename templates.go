package main

const formCreateFile = `<html>

<head>
	<title></title>
	<link rel="stylesheet" href="/css/css.css">
</head>

<body>
	<div class="form">
		<div class="title">
			<p>
				Form
			</p>
		</div>

		<form action="/updatevoc/api/v2/files" enctype="multipart/form-data" method="post">
			<div>
				<label>Description:</label>
				<input type="text" name="description" autocomplete="off">
			</div>
			<div>
				<label>Version:</label>
				<input type="text" name="version" autocomplete="off">
			</div>
			<div>
				<label>Reference:</label>
				<input type="number" name="reference" autocomplete="off">
			</div>
			<div>
				<label>Path:</label>
				<input type="text" name="path" autocomplete="off">
			</div>
			<div>
				<label>Force reboot?:</label>
				<input type="checkbox" name="reboot" value="yes" autocomplete="off">
			</div>
			<div>
				<label>Override?:</label>
				<input type="checkbox" name="override" value="yes">
			</div>
			<div>
				<input type="file" name="fileToUpload" id="fileToUpload">
			</div>
			<div>
				<input type="submit" value="Upload">
			</div>
		</form>
	</div>
</body>

</html>`

const formDeleteFile = `<html>

<head>
    <title></title>
    <link rel="stylesheet" href="/css/css.css">
</head>

<body>
    <div class="form">
        <div class="title">
            <p>
                Show currents files:
            </p>
        </div>
        <form action="/updatevoc/api/v2/files/delete" method="post">
            {{range .}}
                <div>
                    <input type="checkbox" name="files" value="{{.Md5}}">

                    <div class="content">
                        <p><span>Name: </span>{{.Name}}</p>
                        <p><span>Md5: </span>{{.Md5}}</p>
                        <p><span>Devices: </span>
							<ul>
								{{range .DeviceName}}
								<li>{{ .}}</li>
								{{end}}
                            </ul>
                        </p>
                    </div>
                </div>
            {{end}}
            <input type="submit" value="Submit">
        </form>
    </div>
</body>

</html>`

const viewDeviceUpdate = `<html>
<head>
    <title></title>
    <link rel="stylesheet" href="/css/css.css">
</head>

<body class="inter">
    <div class="content">
        <div class="title">
            <p>{{.Name}}</p>
        </div>

		{{range .Updates}}
        <div class="info">
            <div class="title">
                <p>Information</p>

                <p>Date</p>
            </div>

            <div class="conten-info">
                <p>
                    <span>Md5: </span>
                    {{.Md5}}
                </p>
                <p>
                    <span>Filepath: </span>
                    {{.FilePath}}
                </p>
                <p>
                    <span>Description: </span>
                    {{.Description}}
                </p>
                <p>
                    <span>Ref: </span>
                    {{.Ref}}
                </p>
                <p>
                    <span>Version: </span>
                    {{.Version}}
                </p>
                <p>
                    <span>Reboot: </span>
                    {{.Reboot}}
                </p>
                <p>
                    <span>Override: </span>
                    {{.Override}}
                </p>

            </div>
        </div>
    </div>
</body>

</html>`

const listDevices = `<!DOCTYPE html>
	<html>
	<body>
	
	<h1>Show currents devices:</h1>
	
	
	{{range .}}
    <p>
		{{.DeviceName}}
    </p>
    {{range .}}
	
	</body>
</html>`

// const (
// 	formFiledata = `<html>
//     <head>
//     <title></title>
//     </head>
//     <body>
// 		<form action="/updatevoc/api/v2/files" enctype="multipart/form-data" method="post">
// 			<div>
// 				<label>description:</label>
// 				<input type="text" name="description">
// 			</div>
// 			<div>
// 				<label>version:</label>
// 				<input type="text" name="version">
// 			</div>
// 			<div>
// 				<label>reference:</label>
// 				<input type="number" name="reference">
// 			</div>
// 			<div>
// 				<label>path:</label>
// 				<input type="text" name="path">
// 			</div>
// 			<div>
// 				<label>force reboot?:</label>
// 				<input type="checkbox" name="reboot" value="yes">
// 			</div>
// 			<div>
// 				<label>override?:</label>
// 				<input type="checkbox" name="override" value="yes">
// 			</div>
// 			<div>
// 				<input type="file" name="fileToUpload" id="fileToUpload">
// 			</div>
// 			<div>
// 				<input type="submit" value="Upload">
// 			</div>
//         </form>
//     </body>
// </html>`

// 	formDeleteFile = `<!DOCTYPE html>
// 	<html>
// 	<body>

// 	<h1>Show currents files:</h1>

// 	<form action="/updatevoc/api/v2/files/delete" method="post">
// 	{{range .}}

// 		<input type="checkbox" name="files" value="{{.Md5}}">{{.Name}} {{.DeviceName}} {{.Md5}}<br>
// 	{{end}}
// 		<input type="submit" value="Submit">
// 	</form>

// 	</body>
//     </html>`
// )
