#pragma once

#include <parse/parse_node.hpp>
#include <boost/program_options/variables_map.hpp>

namespace po = boost::program_options;

namespace ql::parse {
    class Parser {
    private:
    public:
        std::shared_ptr<ParseNode> parse(po::variables_map& options);

        std::vector<std::string> extractScopes(std::string const& code);

        void recurseScopes(std::string const& code, std::vector<std::string>& scopes);
    };
}
