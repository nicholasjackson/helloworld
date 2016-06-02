require 'minke'
require 'cucumber/rest_api'

discovery = Minke::Docker::ServiceDiscovery.new 'config.yml'
$SERVER_PATH = "http://#{discovery.public_address_for 'helloworld', '8001', :cucumber}"
