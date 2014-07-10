%global debug_package %{nil}

Name:		etcd
Version:	0.4.5
Release:	2%{?dist}
Summary:	A highly-available key value store for shared configuration

License:	ASL 2.0
URL:		https://github.com/coreos/etcd/
Source0:	https://github.com/coreos/%{name}/archive/v%{version}/%{name}-v%{version}.tar.gz
Source1:	etcd.service
Source2:	etcd.socket
Patch1:		etcd-0001-feat-activation-add-socket-activation.patch

BuildRequires:	golang
BuildRequires:	systemd

Requires(post): systemd
Requires(preun): systemd
Requires(postun): systemd

%description
A highly-available key value store for shared configuration.

%prep
%setup -q
sed -i "s/^\(VER=\).*HEAD)/\1%{version}/" ./scripts/release-version
#%patch1 -p1 -b .systemd-activation

%build
export GOPATH=%{gopath}
./build

%install
install -D -p  bin/etcd %{buildroot}%{_bindir}/etcd
#install -t %{buildroot}%{_bindir} etcd
install -D -p -m 0644 %{SOURCE1} %{buildroot}%{_unitdir}/%{name}.service
install -D -p -m 0644 %{SOURCE2} %{buildroot}%{_unitdir}/%{name}.socket

%post
%systemd_post %{name}.service

%preun
%systemd_preun %{name}.service

%postun
%systemd_postun %{name}.service

%files
%{_bindir}/etcd
%{_unitdir}/%{name}.service
%{_unitdir}/%{name}.socket
%doc LICENSE README.md Documentation/internal-protocol-versioning.md

%changelog
* Wed Jul 9 2014 Albert Zhang <zhgwenming@gmail.com> - 0.4.5-2
- rebase etcd to v0.4.5

* Sat Oct 12 2013 Peter Lemenkov <lemenkov@gmail.com> - 0.1.2-2
- Fix typo in the etc.service file

* Sat Oct 12 2013 Peter Lemenkov <lemenkov@gmail.com> - 0.1.2-1
- Ver. 0.1.2
- Integrate with systemd

* Mon Aug 26 2013 Luke Cypret <cypret@fedoraproject.org> - 0.1.1-1
Initial creation
