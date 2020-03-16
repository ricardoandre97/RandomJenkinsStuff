{{ range $i := .Folders -}}
folder('{{ $i.FolderName }}') {
    description('{{ $i.FolderDesc }}')
}
{{ end }}
{{ range $i := .Jobs -}}
{{ if $i.JobFolder }}
pipelineJob('{{ $i.JobFolder }}/{{ $i.JobName }}') {
{{ else }}
pipelineJob('{{ $i.JobName }}') {
{{ end }}
    description('{{ $i.JobDesc }}')
{{ if $i.Params }}
    parameters {
    {{ range $p := $i.Params }}
        {{ $p.Type }}('{{ $p.Name }}', defaultValue = '{{ $p.Value }}', description = '{{ $p.Desc }}')
    {{ end -}}
    }
{{ end }}
    definition {
        cpsScm {
            scm {
                git {
                    remote { 
                        url('{{ $i.GitURL }}')
                        credentials('{{ $i.CredsID }}')
                    }
                    branch('{{ $i.Branch }}')
                    scriptPath('{{ $i.ScriptPath }}')
                    extensions { }
                }
            }
        }
    }
}
{{ end }}