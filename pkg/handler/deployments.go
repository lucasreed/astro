package handler


import (
	log "github.com/sirupsen/logrus"
	appsv1 "k8s.io/api/apps/v1"
  "github.com/reactiveops/dd-manager/conf"
  "text/template"
  "os"
  "bytes"
)


func OnUpdatedDeployment(deployment *appsv1.Deployment) {
  cfg := conf.New()
  monitors := cfg.GetMatchingMonitors(deployment.Annotations, "deployment")
  for _, monitor := range *monitors {
    log.Infof("looping monitor %s", monitor.Name)
    applyDeploymentTemplate(deployment, &monitor)
  }
}


func OnCreatedDeployment(deployment *appsv1.Deployment) {
  cfg := conf.New()
  monitors := cfg.GetMatchingMonitors(deployment.Annotations, "deployment")
  for _, monitor := range *monitors {
    log.Infof("looping monitor %s", monitor.Name)
    applyDeploymentTemplate(deployment, &monitor)
  }
}


func OnDeletedDeployment(deployment *appsv1.Deployment) {
	log.Infof("I finally made it to my deletion handler.")
  cfg := conf.New()
  monitors := cfg.GetMatchingMonitors(deployment.Annotations, "deployment")
  for _, monitor := range *monitors {
    log.Infof("looping monitor %s", monitor.Name)
    applyDeploymentTemplate(deployment, &monitor)
  }
}


func applyDeploymentTemplate(deployment *appsv1.Deployment, monitor *conf.Monitor) {
  var err error
  var tpl bytes.Buffer
  name, _ := template.New("name").Parse(monitor.Name)
  query, _ := template.New("query").Parse(monitor.Query)
  msg, _ := template.New("message").Parse(monitor.Message)
  em, _ := template.New("escalation_message").Parse(monitor.EscalationMessage)

  err = name.Execute(&tpl, deployment)
  if err != nil {
    log.Errorf("Error templating name: %s", err)
  }
  monitor.Name = tpl.String()
  tpl.Reset()

  err = query.Execute(&tpl, deployment)
  if err != nil {
    log.Errorf("Error templating query: %s, err)
  }
  monitor.Query = tpl.String()
  tpl.Reset()

  err = msg.Execute(&tpl, deployment)
  if err != nil {
    log.Error("Error templating message: %s", err)
  }
  monitor.Message = tpl.String()
  tpl.Reset()

  err = em.Execute(&tpl, deployment)
  if err != nil {
    log.Errorf("Error templating escalation message: %s", err)
  }
  monitor.EscalationMessage = tpl.String()
}