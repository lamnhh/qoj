<!DOCTYPE html>
<html>
<head>
    <meta charset="utf-8" />
    <meta http-equiv="X-UA-Compatible" content="IE=edge" />
    <meta name="viewport" content="width=device-width, initial-scale=1" />

    <title>QHH Online Judge</title>

    <link rel="stylesheet" href="/static/index.css">
</head>
<body>
    <form class="submission-form" onsubmit="event.preventDefault(); handleSubmit(this);">
        <label class="submission-form__section">
            <h3>Problem</h3>
            <select name="problemId" required>
                {{range $problem := .}}
                    <option value="{{$problem.Id}}">{{$problem.Code}} - {{$problem.Name}}</option>
                {{end}}
            </select>
        </label>
        <label class="submission-form__section">
            <h3>Source Code</h3>
            <input type="file" name="file">
        </label>
        <button type="submit">Submit</button>
    </form>

    <button type="button" onclick="subscribe();">Subscribe</button>

    <script src="/static/index.js"></script>
</body>
</html>
