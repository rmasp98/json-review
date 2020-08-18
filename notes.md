stuff to check for:
* RBAC
* pod hardening
* PSP use
* Network policy use
* https://kubernetes.io/docs/concepts/security/pod-security-standards/

vulerable clusters:
* https://github.com/raesene/kube_security_lab
* https://www.kubesim.io/
* https://www.bustakube.com/
* https://securekubernetes.com/


## TODO
* Create a predefined querylist
* kubectl pull of config (plugin)
* Create and save query
* query without GUI
* Save vulnXML (plugin)
* Additional fucntions
  * ForValue loop - will insert value into regex of functions?
  * Get where path to it can be defined similar to split
  * FindNodeEquals - return index if key matches key regex and value matches value regex



## Issues
* Secrets in ConfigMap
  * currently only finds private keys but can expand
* Network Policy by namespace
  * Expand to identify namespaces without policy
* Overly Permissive PSP
  * need to expand to cover all parts of PSP 
  * https://kubernetes.io/docs/concepts/policy/pod-security-policy/
* PSP in use?
* Overly permissive network policy
* NodePorts in use?
* Overly permissive Role/ClusterRole? 
  * Any wildcarding maybe
* Certificate authentication? 
  * not sure if we will pick this up here
* RoleBindings to Cluster-admin
* Insecure docker registry 
  * I think you can set this up here so look for that configuraiton without creds and on http
* Kubernetes Auditing 
  * AuditPolicy (https://kubernetes.io/docs/tasks/debug-application-cluster/audit/)
  * Kind AuditSink may give information about where logs sent
  * Write access to this feature gives read access to all cluster data
* resource quotas namespace/pods? 
  * kind=ResourceQuota
* Use of security context 
  * https://kubernetes.io/docs/tasks/configure-pod-container/security-context/
  * allowPrivilegeEscaltion
  * RunAsUser (prevent running as root)
  * Capabilities?
  * Selinux/secomp/apparmour?
* Cloud metadata access 
  * look for network policies not blocking each of the metadata endpoints
* Alpha/Beta features enabled/in use
* Credential rotation
  * Currenrly not sure if collecting this data?
* Port 80 exposed?
