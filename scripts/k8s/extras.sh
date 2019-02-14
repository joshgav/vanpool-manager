# grant cluster-admin rights to kubernetes-dashboard
kubectl create clusterrolebinding \
    kubernetes-dashboard-binding \
    --clusterrole=cluster-admin \
    --serviceaccount=kube-system:kubernetes-dashboard
