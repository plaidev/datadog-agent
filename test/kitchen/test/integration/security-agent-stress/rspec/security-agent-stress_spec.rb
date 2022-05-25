print `cat /etc/os-release`
print `uname -a`

describe 'successfully run stress test against master' do
  it 'displays PASS and returns 0' do
    output = `sudo /tmp/security-agent/stresssuite -test.v -duration 120 1>&2`
    retval = $?
    expect(retval).to eq(0)
  end
end
