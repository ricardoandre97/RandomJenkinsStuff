{{ range $i := . -}}

pipelineJob('{{ $i.JobName }}') {

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