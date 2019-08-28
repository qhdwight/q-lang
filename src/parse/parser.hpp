#pragma once

#include <boost/program_options/variables_map.hpp>

#include <parse/node/parse_node.hpp>

namespace po = boost::program_options;

namespace ql::parse {
    class Parser {
    private:
    public:
        std::shared_ptr<ParseNode> parse(po::variables_map& options);

        std::shared_ptr<ParseNode> getNodes(std::string code);

        void recurseNodes(const std::string& code, const std::weak_ptr<ParseNode>& parent, int depth = 0);
    };
}
