Name:           postfix-aws-cassandra
Version:        1.0.0
Release:        1%{?dist}
Summary:        Postfix Socketmap Daemon for AWS Keyspaces

License:        Apache License
URL:            https://github.com/ethyaan/postfix-aws-cassandra
%global debug_package %{nil}
Source0:        %{name}-%{version}.tar.gz
BuildArch:      x86_64

Requires:       postfix

%description
A Go-based socketmap daemon for Postfix that queries AWS Keyspaces (Cassandra) for mail routing information.

%prep
%setup -q

%build
cd src
GOOS=linux GOARCH=amd64 go build -o postfix-aws-cassandra main.go

%install
install -D -m 0755 %{SOURCE0} %{buildroot}/usr/local/bin/%{name}
install -D -m 0644 postfix-aws-cassandra.service %{buildroot}/usr/lib/systemd/system/postfix-aws-cassandra.service

%files
/usr/local/bin/%{name}
/usr/lib/systemd/system/postfix-aws-cassandra.service

%changelog
* Tue Sep 24 2024 Your Name <you@example.com> - 1.0.0-1
- Initial package