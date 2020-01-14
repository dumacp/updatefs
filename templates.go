package main

const formUpdate = `<html>

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

const formDelete = `<html>

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
    <link rel="stylesheet" href="css.css">
</head>

<body class="inter">
    <div class="content">
        <div class="title">
            <p>{{.Name}}</p>
        </div>

		{{range .}}
        <div class="info">
            <div class="title">
                <p>Information</p>

                <p>Date</p>
            </div>

            <div class="conten-info">
                <p>
                    <span>Md5: </span>
                    b75776f8b5f0786d67908ec9d891af25
                </p>
                <p>
                    <span>Filepath: </span>
                    migracion_test3.zip
                </p>
                <p>
                    <span>Description: </span>
                    test 3
                </p>
                <p>
                    <span>Ref: </span>
                    3
                </p>
                <p>
                    <span>Version: </span>
                    test3
                </p>
                <p>
                    <span>Reboot: </span>
                    true
                </p>
                <p>
                    <span>Override: </span>
                    false
                </p>

            </div>
        </div>

        <div class="info">
            <div class="title">
                <p>Information</p>

                <p>Date</p>
            </div>

            <div class="conten-info">
                <p>
                    <span>Md5: </span>
                    b75776f8b5f0786d67908ec9d891af25
                </p>
                <p>
                    <span>Filepath: </span>
                    migracion_test3.zip
                </p>
                <p>
                    <span>Description: </span>
                    test 3
                </p>
                <p>
                    <span>Ref: </span>
                    3
                </p>
                <p>
                    <span>Version: </span>
                    test3
                </p>
                <p>
                    <span>Reboot: </span>
                    true
                </p>
                <p>
                    <span>Override: </span>
                    false
                </p>

            </div>
        </div>

        <div class="info">
            <div class="title">
                <p>Information</p>

                <p>Date</p>
            </div>

            <div class="conten-info">
                <p>
                    <span>Md5: </span>
                    b75776f8b5f0786d67908ec9d891af25
                </p>
                <p>
                    <span>Filepath: </span>
                    migracion_test3.zip
                </p>
                <p>
                    <span>Description: </span>
                    test 3
                </p>
                <p>
                    <span>Ref: </span>
                    3
                </p>
                <p>
                    <span>Version: </span>
                    test3
                </p>
                <p>
                    <span>Reboot: </span>
                    true
                </p>
                <p>
                    <span>Override: </span>
                    false
                </p>

            </div>
        </div>

        <div class="info">
            <div class="title">
                <p>Information</p>

                <p>Date</p>
            </div>

            <div class="conten-info">
                <p>
                    <span>Md5: </span>
                    b75776f8b5f0786d67908ec9d891af25
                </p>
                <p>
                    <span>Filepath: </span>
                    migracion_test3.zip
                </p>
                <p>
                    <span>Description: </span>
                    test 3
                </p>
                <p>
                    <span>Ref: </span>
                    3
                </p>
                <p>
                    <span>Version: </span>
                    test3
                </p>
                <p>
                    <span>Reboot: </span>
                    true
                </p>
                <p>
                    <span>Override: </span>
                    false
                </p>

            </div>
        </div>
    </div>
</body>

</html>