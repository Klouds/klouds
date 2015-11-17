<div id="welcome">
    <div class="container">
        <div class="title">
            <h1>Your Applications</h1>
            &nbsp;
            {{if .flash.error}}
            <h3>{{.flash.error}}</h3>
            &nbsp;
            {{end}}{{if .flash.notice}}
            <h3>{{.flash.notice}}</h3>
            &nbsp;
            {{end}}
            {{if .Errors}}
            {{range $rec := .Errors}}
            <h3>{{$rec}}</h3>
            {{end}}
            &nbsp;
            {{end}}
        </div>
        <div >
        <h2><a href="/appLaunch" class="button">Launch an App!</a> <br>
        
            <table id="appList">
                <col width="33%">
                <col width="33%">
                <tr>
                    
                   
                        <h2>APPLICATION LIST</h2>
                    
                
                </tr>    

                <tr>
                    <td >Application name<br>
                    -----------------</td>
                    <td>Application link<br>
                        -----------------
                    </td>
                    <td>DELETE APPLICATION<br>
                    ----------------</td>
                </tr>
                {{str2html .AppList}}   
            </table>
        </div>       
    </div>
</div>