#pragma once

#include <boost/program_options/variables_map.hpp>

#include <parser/node/master_node.hpp>

namespace po = boost::program_options;

namespace ql::parser {
    class Parser {
    private:
    public:
        std::shared_ptr<MasterNode> parse(po::variables_map& options);

        std::shared_ptr<MasterNode> getNodes(std::string code);

        void recurseNodes(std::string const& code, std::weak_ptr<AbstractNode> const& parent, int depth = 0);
    };
}
