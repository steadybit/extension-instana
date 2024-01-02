// SPDX-License-Identifier: MIT
// SPDX-FileCopyrightText: 2022 Steadybit GmbH

package extmaintenance

const (
	MaintenanceWindowActionId   = "com.steadybit.extension_instana.maintenance-window"
	maintenanceWindowActionIcon = "data:image/svg+xml;base64,PHN2ZyB3aWR0aD0iMjQiIGhlaWdodD0iMjUiIHZpZXdCb3g9IjAgMCAyNCAyNSIgZmlsbD0ibm9uZSIgeG1sbnM9Imh0dHA6Ly93d3cudzMub3JnLzIwMDAvc3ZnIj48cGF0aCBkPSJNNi4xNyAxNC43MzVjLjY4Ny44MjUgMS45MTIgMS4wNTcgMi44ODYgMS4xNzIuOTIuMTA4IDIuNzgzLjEzNCAyLjc4My4xMzRzMS44NjEtLjAyNSAyLjc4Mi0uMTM0Yy45NzUtLjExNSAyLjE5OC0uMzQ3IDIuODg1LTEuMTcyLjgwNS0uOTY2Ljk5LTIuMjA0IDEuMjIzLTMuMzc0LjM1LTEuNzY2LjM3MS0zLjU4LjA2NC01LjM1NGExLjQxMiAxLjQxMiAwIDAwLS40MzgtLjggMTIuMTYzIDEyLjE2MyAwIDAwLTEuMTQ0LS45MTYgOC41MzQgOC41MzQgMCAwMC0xLjQ0OC0uODY1IDEwLjIwNCAxMC4yMDQgMCAwMC0yLjA3LS43MDNjLS41NTctLjEyLTEuMzQ4LS4yMjMtMS44NTQtLjIyMy0uNTA1IDAtMS4yOTYuMTA0LTEuODUzLjIyMy0uNzE3LjE1NC0xLjQwMi40LTIuMDcuNzAzLS41MTcuMjM0LS45OS41MzYtMS40NDguODY1LS40LjI4Mi0uNzgyLjU4OC0xLjE0NS45MTZhMS40MSAxLjQxIDAgMDAtLjQzOC43OTkgMTQuNjcyIDE0LjY3MiAwIDAwLjA2NSA1LjM1NWMuMjMgMS4xNy40MTUgMi40MDggMS4yMiAzLjM3NHptOC44NzItMS42ODJjLjA0NS0uNTg3LjQ1Ni0xLjAzOC45MTgtMS4wMDkuNDYxLjAzLjguNTI5Ljc1NCAxLjExNS0uMDQ0LjU4Ny0uNDU1IDEuMDM4LS45MTYgMS4wMDktLjQ2Mi0uMDMtLjgtLjUzLS43NTYtMS4xMTV6bS03LjMxOS0xLjAwOWMuNDYyLS4wMzIuODcuNDE3LjkxIDEuMDAzLjA0MS41ODYtLjMgMS4wODgtLjc2MiAxLjEyLS40NjEuMDMzLS44NjktLjQxNi0uOTEtMS4wMDItLjA0LS41ODcuMzAxLTEuMDg4Ljc2Mi0xLjEyem0xMi42OTItLjc0NGwtLjA5LS4wMThjLjAzNy0uMzcxLjA1LS43NDQuMDQyLTEuMTE3LS4wMTItLjM5LS4xMzItMi4wMTctLjQ1Ny0yLjk3Ni0uMTYyLS40NzctLjMzNi0uOTM0LS42NTctMS4zNDYtLjAzNC0uMDQ0LS4wNzItLjA5LS4xMS0uMTM3YS4wNjEuMDYxIDAgMDAtLjEwOS4wNTNjLjQxNSAxLjc4OS40IDMuNzg0LjEwNSA1LjU2NC0uMTkyIDEuMTU5LS40NiAyLjUxMi0xLjA3IDMuNTA1LS42NzEgMS4wOTctMS45MDkgMS4zNTQtMy4wMjIgMS41MjUtMS4wNTguMTYyLTMuMjEuMTg2LTMuMjEuMTg2cy0yLjE1Mi0uMDI0LTMuMjEtLjE4NmMtMS4xMTItLjE3MS0yLjM1LS40MjgtMy4wMjItMS41MjYtLjYwOC0uOTk0LS44NzgtMi4zNDktMS4wNy0zLjUwNS0uMjkzLTEuNzgtLjMwOS0zLjc3NC4xMDYtNS41NjVhLjA2MS4wNjEgMCAwMC0uMTA5LS4wNTNjLS4wNC4wNDgtLjA3Ni4wOTMtLjExLjEzOC0uMzIuNDExLS40OTUuODY3LS42NTcgMS4zNDYtLjMyNS45NTgtLjQ0NSAyLjU4NS0uNDU3IDIuOTc2LS4wMDguMzczLjAwNi43NDUuMDQxIDEuMTE3bC0uMDkuMDE4Yy0uMTY4LjAzNi0uMjguMTc0LS4yNTYuMzIybC41MzkgMy40MjNjLjAyMy4xNDguMTcyLjI1Ny4zNDYuMjUzbC4zOS0uMDA5Yy4wODIuMTkuMTczLjM3Ni4yNzUuNTU3LjI0Mi40MzQuNTkuNzU1IDEuMDEyIDEuMDA1LjQwNS4yNDEuODUuMzcgMS4zMDUuNDczLjUzMS4xMiAxLjA3LjE5MiAxLjYxLjI1M2wuNTMyLjA2NWMuMDA3IDAgLjAxNC4wMDQuMDIuMDFhLjAzMy4wMzMgMCAwMS4wMDUuMDQuMDM0LjAzNCAwIDAxLS4wMTcuMDE1Yy0uNDIuMTIzLTEuMzIxLjUzOC0xLjcxNC45MWE1Ljg4NiA1Ljg4NiAwIDAwLS45NjIgMS4wNjNjLS4yMzYuMzQxLS40NDcuNjk5LS41NTEgMS4xMDV2LjAwN2EuNjkuNjkgMCAwMC40NTcuODE1YzEuNzEzLjU3NSAzLjYwMy44OTQgNS41ODkuODk0IDEuOTg2IDAgMy44NzUtLjMxOSA1LjU4OC0uODk0YS42OS42OSAwIDAwLjQ1OC0uODE2bC0uMDAxLS4wMDZjLS4xMDQtLjQwNi0uMzE1LS43NjQtLjU1MS0xLjEwNWE1Ljg4NCA1Ljg4NCAwIDAwLS45NjUtMS4wNThjLS4zOTMtLjM3Mi0xLjI5My0uNzg4LTEuNzE0LS45MTFhLjAzNS4wMzUgMCAwMS0uMDE3LS4wMTQuMDM0LjAzNCAwIDAxLjAyNS0uMDVjLjE0OS0uMDIuMzktLjA0OS41MzEtLjA2Ni41NDItLjA2MyAxLjA4LS4xMzQgMS42MTEtLjI1Mi40NTUtLjEwMy45LS4yMzMgMS4zMDYtLjQ3NC40MjItLjI1Ljc3LS41NzIgMS4wMTEtMS4wMDUuMTAyLS4xODEuMTk0LS4zNjcuMjc2LS41NTdsLjM5LjAxYy4xNzIuMDA0LjMyMi0uMTA1LjM0NS0uMjUzbC41MzktMy40MjRjLjAyNC0uMTUtLjA4Ny0uMjktLjI1Ni0uMzI1eiIgZmlsbD0iY3VycmVudENvbG9yIi8+PC9zdmc+"
)
