package utils

import (
    "context"
    "log"

    corev1 "k8s.io/api/core/v1"
    metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
    "k8s.io/client-go/kubernetes"
    "k8s.io/client-go/tools/clientcmd"
)

func GetSvcAcctToken(kubeConfigPath string, namespace string, svcAcctName string) (*string, error) {
    secret, err := GetSvcAcctSecret(kubeConfigPath, namespace, svcAcctName+"-token")
    if err != nil {
        log.Printf("Problem with GetSvcAcctSecret: %v", err)
        return nil, err
    }
    token := string(secret.Data["token"])

    return &token, err
}

func GetSvcAcctSecret(kubeConfigPath string, namespace string, secretName string) (*corev1.Secret, error) {
    config, err := clientcmd.BuildConfigFromFlags("", kubeConfigPath)
    if err != nil {
        log.Fatalf("Problem with BuildConfigFromFlags: %v", err)
        return nil, err
    }

    k8sClient, err := kubernetes.NewForConfig(config)
    if err != nil {
        log.Fatalf("Problem with NewForConfig: %v", err)
        return nil, err
    }

    secret, err := k8sClient.CoreV1().Secrets(namespace).Get(
        context.Background(),
        secretName,
        metav1.GetOptions{},
    )

    return secret, err
}
